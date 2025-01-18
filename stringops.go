package main

import (
	"math/rand"
	"os"
	"slices"
	"strings"
)

var zalgUp = []rune{
	'̍', '̎', '̄', '̅',
	'̿', '̑', '̆', '̐',
	'͒', '͗', '͑', '̇',
	'̈', '̊', '͂', '̓',
	'̈', '͊', '͋', '͌',
	'̃', '̂', '̌', '͐',
	'̀', '́', '̋', '̏',
	'̒', '̓', '̔', '̽',
	'̉', 'ͣ', 'ͤ', 'ͥ',
	'ͦ', 'ͧ', 'ͨ', 'ͩ',
	'ͪ', 'ͫ', 'ͬ', 'ͭ',
	'ͮ', 'ͯ', '̾', '͛',
	'͆', '̚',
}

var zalgDown = []rune{
	'̖', '̗', '̘', '̙',
	'̜', '̝', '̞', '̟',
	'̠', '̤', '̥', '̦',
	'̩', '̪', '̫', '̬',
	'̭', '̮', '̯', '̰',
	'̱', '̲', '̳', '̹',
	'̺', '̻', '̼', 'ͅ',
	'͇', '͈', '͉', '͍',
	'͎', '͓', '͔', '͕',
	'͖', '͙', '͚', '̣',
}

var zalgMid = []rune{
	'̕', '̛', '̀', '́',
	'͘', '̡', '̢', '̧',
	'̨', '̴', '̵', '̶',
	'͜', '͝', '͞',
	'͟', '͠', '͢', '̸',
	'̷', '͡', '҉',
}

// sort sorts a list alphabetically
func sort() string {
	split := strings.Split(os.Getenv("ESPANSO_CLIPBOARD"), "\n")
	slices.Sort(split)
	return strings.Join(split, "\n")
}

func capitalise() string {
	return strings.ToUpper(os.Getenv("ESPANSO_CLIPBOARD"))
}

func uncaptialise() string {
	return strings.ToLower(os.Getenv("ESPANSO_CLIPBOARD"))
}

func randRune(arr []rune) rune {
	return arr[rand.Intn(len(arr))]
}

func zalgo() string {
	input := os.Getenv("ESPANSO_CLIPBOARD")
	output := ""
	for _, char := range input {
		output += string(randRune(zalgUp))
		output += string(randRune(zalgMid))
		output += string(char)
		output += string(randRune(zalgDown))
	}
	return output
}

var regular = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
var transformations = map[string][]rune{
	"fullwidth":   []rune("ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ"),
	"bold":        []rune("𝐀𝐁𝐂𝐃𝐄𝐅𝐆𝐇𝐈𝐉𝐊𝐋𝐌𝐍𝐎𝐏𝐐𝐑𝐒𝐓𝐔𝐕𝐖𝐗𝐘𝐙𝐚𝐛𝐜𝐝𝐞𝐟𝐠𝐡𝐢𝐣𝐤𝐥𝐦𝐧𝐨𝐩𝐪𝐫𝐬𝐭𝐮𝐯𝐰𝐱𝐲𝐳"),
	"italic":      []rune("𝘈𝘉𝘊𝘋𝘌𝘍𝘎𝘏𝘐𝘑𝘒𝘓𝘔𝘕𝘖𝘗𝘘𝘙𝘚𝘛𝘜𝘝𝘞𝘟𝘠𝘡𝘢𝘣𝘤𝘥𝘦𝘧𝘨𝘩𝘪𝘫𝘬𝘭𝘮𝘯𝘰𝘱𝘲𝘳𝘴𝘵𝘶𝘷𝘸𝘹𝘺𝘻"),
	"pirate":      []rune("𝕬𝕭𝕮𝕯𝕰𝕱𝕲𝕳𝕴𝕵𝕶𝕷𝕸𝕹𝕺𝕻𝕼𝕽𝕾𝕿𝖀𝖁𝖂𝖃𝖄𝖅𝖆𝖇𝖈𝖉𝖊𝖋𝖌𝖍𝖎𝖏𝖐𝖑𝖒𝖓𝖔𝖕𝖖𝖗𝖘𝖙𝖚𝖛𝖜𝖝𝖞𝖟"),
	"fancy":       []rune("𝓐𝓑𝓒𝓓𝓔𝓕𝓖𝓗𝓘𝓙𝓚𝓛𝓜𝓝𝓞𝓟𝓠𝓡𝓢𝓣𝓤𝓥𝓦𝓧𝓨𝓩𝓪𝓫𝓬𝓭𝓮𝓯𝓰𝓱𝓲𝓳𝓴𝓵𝓶𝓷𝓸𝓹𝓺𝓻𝓼𝓽𝓾𝓿𝔀𝔁𝔂𝔃"),
	"smallcaps":   []rune("ᴀʙᴄᴅᴇꜰɢʜɪᴊᴋʟᴍɴᴏᴩQʀꜱᴛᴜᴠᴡxYᴢᴀʙᴄᴅᴇꜰɢʜɪᴊᴋʟᴍɴᴏᴩqʀꜱᴛᴜᴠᴡxyᴢ"),
	"subscript":   []rune("ᴬᴮᶜᴰᴱᶠᴳᴴᴵᴶᴷᴸᴹᴺᴼᴾQᴿˢᵀᵁⱽᵂˣʸᶻᵃᵇᶜᵈᵉᶠᵍʰⁱʲᵏˡᵐⁿᵒᵖqʳˢᵗᵘᵛʷˣʸᶻ"),
	"superscript": []rune("ₐBCDₑFGₕᵢⱼₖₗₘₙₒₚQᵣₛₜᵤᵥWₓYZₐbcdₑfgₕᵢⱼₖₗₘₙₒₚqᵣₛₜᵤᵥwₓyz"),
}

func transform(outputType string) func() string {
	return func() string {
		input := os.Getenv("ESPANSO_CLIPBOARD")
		output := ""
		for _, char := range input {
			index := strings.Index(regular, string(char))
			if index == -1 {
				output += string(char)
			} else {
				output += string(transformations[outputType][index])
			}
		}
		return output
	}
}
