package app

import "math/rand"

func chickenSoup() string {
	var soups = []string{
		"百分之八十的成功来自于出席。",
		"障碍与失败，是通往成功最稳靠的踏脚石，肯研究利用它们，便能从失败中培养出成功。",
		"奋斗没有终点，任何时候都是一个起点。",
		"每一次努力都会离梦想更近一步。",
		"永不言败，是成功者的最佳品格。",
		"一个人几乎可以在任何他怀有无限热忱的事情上成功。",
		"Hit the ball hard and good things happen.",
		"Figure out your opponent, “hack your opponent”, their tendencies, and how to exploit them.",
		"If you understand why you lost, it’s a win in itself.",
		"Back to basics. Always go back to the fundamentals.",
		"Everybody has talent, but ability takes hard work. -- Michael Jordan",
	}

	return soups[rand.Intn(len(soups))]
}
