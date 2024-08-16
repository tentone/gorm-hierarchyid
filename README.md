# Gorm hierarchyid
 - Library to handle hierarchyid type in SQL Server and go.
   - Generation and parsing of hierarchyid type in go.
   - Type wrapper for usage with gorm ORM.
 - The [hierarchyid](https://learn.microsoft.com/en-us/sql/relational-databases/hierarchical-data-sql-server?view=sql-server-ver16) is data to represent a position in a hierarchy in SQL Server.
   - It is a variable length type with reduced storage requirements.
 - Encodes the position in the hierarchy as a list of indexes
   - For example in the tree below the path to `E` is `/1/1/2/`
   - Indexes can be used to sort elements inside of a tree level.
 
<img src="./readme/tree.png" width="600"/>

## How it works
 - The `HierarchyID` is defined as a `[]int64` in go.
 - When serialized into JSON a textual representation is used for readability.
   - Represented as list separated by `/`. (e.g. `/1/2/3/4/5/`)
 - Each element in the slice represents a level in the hierarchy.
 - An empty slice represents the root of the hierarchy.
   - Elements placed in the root should not use an empty list.
   - They should instead by represented by `/1/`, `/2/`, etc.

## Installation
 - The library can be installed using `go get`.
 - 
  ```bash
  go get github.com/tentone/gorm-hierarchyid
  ```

## Model definition
 - Declare `HierarchyID` type in your gorm model, there is no need to specify the DB data type.
 - Is is recommended to also mark the field as `unique` to avoid duplicates.
 - The library will handle the serialization and deserialization of the field to match the SQL Server `hierarchyid` type.
    ```go
    type Model struct {
        gorm.Model

        Path HierarchyID `gorm:"unique;not null;"`
    }
    ```
 - In some scenarios it might be usefull to also keep a tradicional relationship to the parent.
   - This can be done by adding a `ParentID` field to the model.
   - It ensures that the tree is consistent and that actions (e.g. delete) are cascaded to the children.
   - Some operations might also be easier to perform with the parent relationship.
    ```go
    type Model struct {
      gorm.Model

      Path HierarchyId `gorm:"unique;not null;"`

      ParentID uint              `gorm:"index"`
      Parent   *TestParentsTable `foreignKey:"parent_id;references:id;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE;"`
    }
    ```

## Usage

### Create
 - Elements can be added to the tree as regular entries
 - Just make sure that the tree indexes are filled correctly, indexes dont need to be sequential.
  ```go
  db.Create(&Table{Path: HierarchyID{Data: []int64{1}}})
  db.Create(&Table{Path: HierarchyID{Data: []int64{1, 1}}})
  db.Create(&Table{Path: HierarchyID{Data: []int64{1, 1, 2}}})
  ```

### Get Ancestors
 - To get all parents of a node use the `GetAncestors` method.
 - The method will return a slice with all the parents of the node. This can be used as param for a query.
  ```go
  db.Model(&Table{}).Where("[path] IN (?)", child.Path.GetAncestors()).Find(&parents)
  ```
 - Its also possible to get parents with the SQL version of the [`GetAncestor`](https://learn.microsoft.com/en-us/sql/t-sql/data-types/getancestor-database-engine?view=sql-server-ver16) method.
 - Example on getting the parent of an element.
  ```go
  db.Model(&Table{}).Where("[path] = ?.GetAncestor(1)", child.Path).Find(&parent)
  ```

### Get Descendants
 - To get all children of a node use the [`IsDescendantOf`](https://learn.microsoft.com/en-us/sql/t-sql/data-types/isdescendantof-database-engine?view=sql-server-ver16) method in SQL.
 - Example on getting all children of a node (including the node itself).
  ```go
  elements := []Table{}
  db.Where("[path].IsDescendantOf(?)=1", HierarchyId{Data: []int64{1, 2}}).Find(&elements)
  ```
 - It is also possible to filter the children based on sub-levels.
 - Example on getting all nodes from root where at least one of the sub-level has a name that contains the text 'de'
  ```sql
  SELECT *
  FROM "table" as a
  WHERE ([path].GetLevel()=1 AND [path].IsDescendantOf('/')=1) AND
  (SELECT COUNT(*) FROM "table" AS b WHERE b.path.IsDescendantOf(a.path)=1 AND b.name LIKE '%de%')>0
  ```
 - The `GetLevel` method can be used to filter nodes based on their level in the hierarchy. Also available in SQL with the same name [`GetLevel`](https://learn.microsoft.com/en-us/sql/t-sql/data-types/getlevel-database-engine?view=sql-server-ver16).
 - A more generic version of the same code presented above writen in go.

  ```go
  root := GetRoot()
  subQuery := db.Table("table AS b").Select("COUNT(*)").Where("[b].[path].IsDescendantOf([a].[path])=1 AND [b].[name] LIKE '%de%'")
  conn = db.Table("table AS a").
    Where("[a].[path].GetLevel()=? AND [a].[path].IsDescendantOf(?)=1 AND (?)>0", root.GetLevel()+1, root, subQuery).
    Find(&elements)
  ```


### Move nodes
 - To move a node to a new parent there is the `GetReparentedValue` method that receives the old parent and new parent and calculates the new hierarchyid value.
 - Example on moving a node to a new parent.
  ```go
  db.Model(&Table{}).Where("[id] = ?", id).Update("[path]=?", node.Path.GetReparentedValue(oldParent.Path, newParent.Path))
  ```

## Resources
 - [adamil.net - How the SQL Server hierarchyid data type works (kind of)](http://www.adammil.net/blog/v100_how_the_SQL_Server_hierarchyid_data_type_works_kind_of_.html)
 - [hierarchyid data type method reference](https://learn.microsoft.com/en-us/sql/t-sql/data-types/hierarchyid-data-type-method-reference?view=sql-server-ver16&redirectedfrom=MSDN)
  - .NET Implementation ([Logic](https://github.com/dotMorten/Microsoft.SqlServer.Types/tree/main/src/Microsoft.SqlServer.Types/SqlHierarchy) + [Interface](https://github.com/dotMorten/Microsoft.SqlServer.Types/blob/main/src/Microsoft.SqlServer.Types/SqlHierarchyId.cs))

## License
 - The project is distributed using a MIT license. Available on the project repository.
