package slicex

// 分批次执行
func BatchExec(start, end, limit int, f func(int, int) error) error {
	if start < end {
		for start < end {
			if start+limit > end {
				limit = end - start
			}
			if err := f(start, start+limit); err != nil {
				return err
			}
			start += limit
		}
	}
	return nil
}
