//go:build !solution

package hogwarts

func GetCourseList(prereqs map[string][]string) (res []string) {
	learned := make(map[string]int)
	for course := range prereqs {
		if v := learned[course]; v == 2 {
			continue
		}
		res = append(res, dfs(course, prereqs, learned)...)
	}
	return
}

func dfs(course string, prereqs map[string][]string, learned map[string]int) (res []string) {
	learned[course] = 1
	for _, preqCourse := range prereqs[course] {
		v, ok := learned[preqCourse]
		if ok && v == 2 {
			continue
		} else if ok && v == 1 {
			panic("cycle")
		}
		res = append(res, dfs(preqCourse, prereqs, learned)...)
	}
	res = append(res, course)
	learned[course] = 2
	return
}
