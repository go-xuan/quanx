package sqlx

import (
	"fmt"
	"testing"
)

func TestSqlFmt(t *testing.T) {
	var sql = `SELECT DISTINCT '123', '456', '789', good(t1.aaa, ',') as aaa, t1.bbb as bbb, good(t1.ccc, ',') as ccc 
FROM t_abc_2 as t1 left join t_abc_2 t2 on t1.aaa = t2.aaa right join (select aaa,bbb,ccc from t_abc_3 where aaa = '1') t2 on t2.aaa = t3.aaa 
where t1.aaa = '1'and (t2.bbb = '1' or t3.ccc = '1') group by t1.aaa, t2.bbb, t3.ccc having t1.aaa > 1 and (t2.bbb <= 1 or t3.ccc = 0)
order by t1.aaa ASC, t2.bbb DESC, t3.ccc ASC limit 10 offset 20`
	fmt.Println(Format(sql).String())
}
