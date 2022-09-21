package accesskit

import (
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

// 补充 gorm/gen 辅助

// 合并多个查询表达式，转成 sql 与 参数对象
//  columns: field.Expr, string, clause.Expr
func buildExpr(stmt *gorm.Statement, columns ...interface{}) (query []string, args []interface{}) {
	for _, e := range columns {
		switch v := e.(type) {
		case field.Expr:
			sql, vars := v.BuildWithArgs(stmt)
			query = append(query, sql.String())
			args = append(args, vars...)
		case string:
			query = append(query, v)
		case clause.Expr:
			query = append(query, v.SQL)
			args = append(args, v.Vars...)
		}
	}
	return query, args
}

// Select 字段
//  columns: field.Expr, string, clause.Expr
func Select(do *gen.DO, columns ...interface{}) {
	db := do.UnderlyingDB()
	query, args := buildExpr(db.Statement, columns...)
	db = db.Select(strings.Join(query, ","), args...)
	do.ReplaceDB(db)
}

// SelectAppend 增加额外的 Select 字段
//  columns: 支持类型 field.Expr, string, clause.Expr
//
//  示例:
//
//     accesskit.SelectAppend(d2.(*gen.DO), clause.Expr{SQL: "ifnull(sum(cnt1)-sum(cnt2),0) cnt1"})
//
//     accesskit.SelectAppend(&dao1.DO,
//       clause.Expr{
//         SQL: "ifnull(?,?) " + dbPMS.Status.ColumnName().String(),
//         Vars: []interface{}{dbPM.Status.RawExpr(), dbPMS.Status.RawExpr()},
//       },
// 	     clause.Expr{
//	       SQL: "max(if(?=?,id,0)) id",
//	       Vars: []interface{}{dbVR.CheckinID.RawExpr(), checkinID},
//	     },
//     }
//
func SelectAppend(do *gen.DO, columns ...interface{}) {
	db := do.UnderlyingDB()
	query, args := buildExpr(db.Statement, columns...)
	if c1, ok := db.Statement.Clauses["SELECT"]; ok && c1.Expression != nil {
		switch v := c1.Expression.(type) {
		case clause.Expr:
			query = append([]string{v.SQL}, query...)
			args = append(v.Vars, args...)
		case clause.NamedExpr:
			query = append([]string{v.SQL}, query...)
			args = append(v.Vars, args...)
		default:
			// NOTE:
		}
	} else {
		query = append(db.Statement.Selects, query...)
	}
	if do.TableName() != "" && db.Statement.TableExpr == nil {
		db = db.Table(do.TableName())
	}
	db = db.Select(strings.Join(query, ","), args...)
	do.ReplaceDB(db)
}

func ColsNamesByExpr(expr ...field.Expr) []string {
	names := make([]string, len(expr))
	for i := range expr {
		names[i] = string(expr[i].ColumnName())
	}
	return names
}
