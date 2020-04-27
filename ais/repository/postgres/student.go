package postgres

// CREATE TABLE Студент(
//     id SERIAL,
//     id_группы int NOT NULL references Группа(id) ON DELETE CASCADE,
//     имя varchar(50) NOT NULL,
//     фамилия varchar(50) NOT NULL,
//     отчество varchar(50),
//     CONSTRAINT студент_pk PRIMARY KEY (id)
// );

type Student struct {
}
