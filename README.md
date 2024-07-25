# Gorm hierarchyid
 - The hierarchyid is data to represent a position in a hierarchy.
 - It is a variable length type with reduced storage requirements.
 - Handle hierarchyid type in SQL Server and go.
 - Implements generation and parsing of hierarchyid type in go.
 - Implements a type wrapper for usage with gorm.

## Usage
 - 
 - 

## Resources
 - [adamil.net - How the SQL Server hierarchyid data type works (kind of)](http://www.adammil.net/blog/v100_how_the_SQL_Server_hierarchyid_data_type_works_kind_of_.html)
 - [hierarchyid data type method reference](https://learn.microsoft.com/en-us/sql/t-sql/data-types/hierarchyid-data-type-method-reference?view=sql-server-ver16&redirectedfrom=MSDN)
  - .NET Implementation ([Logic](https://github.com/dotMorten/Microsoft.SqlServer.Types/tree/main/src/Microsoft.SqlServer.Types/SqlHierarchy) + [Interface](https://github.com/dotMorten/Microsoft.SqlServer.Types/blob/main/src/Microsoft.SqlServer.Types/SqlHierarchyId.cs))
 - [Github](https://github.com/dotMorten/Microsoft.SqlServer.Types/pull/33) pull request Long type support

## License
 - The project is distributed using a MIT license. Available on the project repository.
