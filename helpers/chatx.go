package helpers

import "math/rand"

var userNameColors = []string{
	"#00FFFF",
	"#F0FFFF",
	"#89CFF0",
	"#7393B3",
	"#088F8F",
	"#0096FF",
	"#5F9EA0",
	"#0047AB",
	"#6495ED",
	"#00FFFF",
	"#00008B",
	"#6F8FAF",
	"#1434A4",
	"#7DF9FF",
	"#6082B6",
	"#00A36C",
	"#3F00FF",
	"#5D3FD3",
	"#ADD8E6",
	"#191970",
	"#000080",
	"#1F51FF",
	"#A7C7E7",
	"#B6D0E2",
	"#96DED1",
	"#4169E1",
	"#0F52BA",
	"#9FE2BF",
	"#87CEEB",
	"#4682B4",
	"#008080",
	"#40E0D0",
	"#0437F2",
	"#40B5AD",
	"#0818A8",
}

func GetRandomChatColor() string {
	ix := rand.Intn(len(userNameColors))
	return userNameColors[ix]
}
