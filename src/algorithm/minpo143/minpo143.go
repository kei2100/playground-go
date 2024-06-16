package minpo143

import (
	"fmt"
	"time"
)

// Minpo143 は、民法143条に基づき fromDate から toDate までの経過月数を返します。
// Date は `YYYY-MM-DD` 形式である必要があります（time.DateOnly）
//
// 月の初日から起算する場合、満了日は最終月の末日になります
//   - 例1: 1月1日から起算して2か月は、平年なら2月28日、閏年なら2月29日が満了日
//   - 例2: 1月1日から起算して3か月は、3月31日が満了日
//
// 月の途中から起算し、最終月に応当日がある場合、満了日は最終月の応当日の前日になります
//   - 例1: 1月20日から起算して2か月は、3月19日が満了日
//   - 例2: 1月31日から起算して2か月は、3月30日が満了日
//
// 月の途中から起算し、最終月に応当日がない場合、満了日は最終月の末日になります
//   - 例1: 1月31日から起算して1か月は、平年なら2月28日、閏年なら2月29日が満了日
//   - 例2: 3月31日から起算して1か月は、4月30日が満了日
func Minpo143(fromDate, toDate string) int {
	from, err := time.Parse(time.DateOnly, fromDate)
	if err != nil {
		panic(fmt.Sprintf("minpo143: invalid from date: %s", fromDate))
	}
	to, err := time.Parse(time.DateOnly, toDate)
	if err != nil {
		panic(fmt.Sprintf("minpo143: invalid to date: %s", toDate))
	}
	if from.After(to) {
		return 0
	}
	fy, fm, fd := from.Year(), from.Month(), from.Day()
	to = to.AddDate(0, 0, 1)
	ty, tm, td := to.Year(), to.Month(), to.Day()
	sumMonth := (ty-fy)*12 + int(tm) - int(fm)
	if td < fd {
		sumMonth--
	}
	return sumMonth
}
