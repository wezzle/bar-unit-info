package gamedata

import "github.com/wezzle/bar-unit-info/gamedata/types"

func GetLabGrid() types.LabGrid {
	return labGridData
}

var labGridData types.LabGrid = types.LabGrid{"armaap":types.GridRow{types.GridCol{"armaca", "armhawk", "armbrawl", "armpnix"}, types.GridCol{"armawac", "armdfly", "armlance", "armsfig2"}, types.GridCol{"armliche", "armblade", "armstil", ""}}, "armalab":types.GridRow{types.GridCol{"armack", "armfark", "armfast", "armspy"}, types.GridCol{"armmark", "armaser", "armzeus", "armmav"}, types.GridCol{"armfido", "armsnipe", "armaak", "armfboy"}}, "armamsub":types.GridRow{types.GridCol{"armbeaver", "armdecom", "armpincer", ""}, types.GridCol{"armcroc", "", "", ""}, types.GridCol{"", "armjeth", "armaak", ""}}, "armap":types.GridRow{types.GridCol{"armca", "armfig", "armkam", "armthund"}, types.GridCol{"armpeep", "armatlas", "", ""}, types.GridCol{"", "", "", ""}}, "armasy":types.GridRow{types.GridCol{"armacsub", "armmls", "armcrus", "armmship"}, types.GridCol{"armcarry", "armsjam", "armbats", "armepoch"}, types.GridCol{"armsubk", "armserp", "armaas", ""}}, "armavp":types.GridRow{types.GridCol{"armacv", "armconsul", "armbull", "armmart"}, types.GridCol{"armseer", "armjam", "armmanni", "armst"}, types.GridCol{"armlatnk", "armcroc", "armyork", "armmerl"}}, "armfhp":types.GridRow{types.GridCol{"armch", "", "armsh", ""}, types.GridCol{"armanac", "armmh", "", ""}, types.GridCol{"", "", "armah", ""}}, "armhp":types.GridRow{types.GridCol{"armch", "", "armsh", ""}, types.GridCol{"armanac", "armmh", "", ""}, types.GridCol{"", "", "armah", ""}}, "armlab":types.GridRow{types.GridCol{"armck", "armrectr", "armpw", "armflea"}, types.GridCol{"armrock", "armham", "armwar", ""}, types.GridCol{"", "", "armjeth", ""}}, "armplat":types.GridRow{types.GridCol{"armcsa", "armsfig", "armsaber", "armsb"}, types.GridCol{"armsehak", "armseap", "", ""}, types.GridCol{"", "", "", ""}}, "armshltx":types.GridRow{types.GridCol{"armmar", "armraz", "armvang", "armthor"}, types.GridCol{"armbanth", "armlun", "", ""}, types.GridCol{"", "", "", ""}}, "armsy":types.GridRow{types.GridCol{"armcs", "armrecl", "armdecade", ""}, types.GridCol{"armpship", "armroy", "", ""}, types.GridCol{"armsub", "", "armpt", ""}}, "armvp":types.GridRow{types.GridCol{"armcv", "armmlv", "armflash", "armfav"}, types.GridCol{"armstump", "armjanus", "armart", ""}, types.GridCol{"armbeaver", "armpincer", "armsam", ""}}, "coraap":types.GridRow{types.GridCol{"coraca", "corvamp", "corape", "corhurc"}, types.GridCol{"corawac", "corseah", "cortitan", "corsfig2"}, types.GridCol{"corcrw", "corcrwh", "", ""}}, "coralab":types.GridRow{types.GridCol{"corack", "corfast", "corpyro", "corspy"}, types.GridCol{"corvoyr", "corspec", "corcan", "corhrk"}, types.GridCol{"cormort", "corroach", "coraak", "corsumo"}}, "coramsub":types.GridRow{types.GridCol{"cormuskrat", "cordecom", "corgarp", ""}, types.GridCol{"corseal", "corparrow", "", ""}, types.GridCol{"", "corcrash", "coraak", ""}}, "corap":types.GridRow{types.GridCol{"corca", "corveng", "corbw", "corshad"}, types.GridCol{"corfink", "corvalk", "", ""}, types.GridCol{"", "", "", ""}}, "corasy":types.GridRow{types.GridCol{"coracsub", "cormls", "corcrus", "cormship"}, types.GridCol{"corcarry", "corsjam", "corbats", "corblackhy"}, types.GridCol{"corshark", "corssub", "corarch", ""}}, "coravp":types.GridRow{types.GridCol{"coracv", "corban", "correap", "cormart"}, types.GridCol{"corvrad", "coreter", "corgol", "cortrem"}, types.GridCol{"corseal", "corparrow", "corsent", "corvroc"}}, "corfhp":types.GridRow{types.GridCol{"corch", "", "corsh", ""}, types.GridCol{"corsnap", "cormh", "corhal", ""}, types.GridCol{"", "", "corah", ""}}, "corgant":types.GridRow{types.GridCol{"corcat", "corkarg", "corshiva", "corkorg"}, types.GridCol{"corjugg", "corsok", "", ""}, types.GridCol{"", "", "", ""}}, "corhp":types.GridRow{types.GridCol{"corch", "", "corsh", ""}, types.GridCol{"corsnap", "cormh", "corhal", ""}, types.GridCol{"", "", "corah", ""}}, "corlab":types.GridRow{types.GridCol{"corck", "cornecro", "corak", ""}, types.GridCol{"corstorm", "corthud", "", ""}, types.GridCol{"", "", "corcrash", ""}}, "corplat":types.GridRow{types.GridCol{"corcsa", "corsfig", "corcut", "corsb"}, types.GridCol{"corhunt", "corseap", "", ""}, types.GridCol{"", "", "", ""}}, "corsy":types.GridRow{types.GridCol{"corcs", "correcl", "coresupp", ""}, types.GridCol{"corpship", "corroy", "", ""}, types.GridCol{"corsub", "", "corpt", ""}}, "corvp":types.GridRow{types.GridCol{"corcv", "cormlv", "corgator", "corfav"}, types.GridCol{"corraid", "corlevlr", "corwolv", ""}, types.GridCol{"cormuskrat", "corgarp", "cormist", ""}}, "legaap":types.GridRow{types.GridCol{"legaca", "legionnaire", "legvenator", ""}, types.GridCol{"legmineb", "legnap", "legphoenix", "cortitan"}, types.GridCol{"legfort", "legstronghold", "legwhisper", ""}}, "legalab":types.GridRow{types.GridCol{"legack", "legaceb", "legstr", "corspy"}, types.GridCol{"corvoyr", "corspec", "leginfestor", "legsrail"}, types.GridCol{"legbart", "corroach", "legshot", "leginc"}}, "legamsub":types.GridRow{types.GridCol{"legotter", "legdecom", "legamphtank", ""}, types.GridCol{"", "legfloat", "", ""}, types.GridCol{"", "corcrash", "coraak", ""}}, "legap":types.GridRow{types.GridCol{"legca", "legfig", "legmos", "legkam"}, types.GridCol{"legcib", "legatrans", "", ""}, types.GridCol{"", "", "", ""}}, "legavp":types.GridRow{types.GridCol{"legacv", "legmrv", "legaskirmtank", "legamcluster"}, types.GridCol{"corvrad", "coreter", "legaheattank", "leginf"}, types.GridCol{"legfloat", "legmed", "corsent", "legavroc"}}, "legfhp":types.GridRow{types.GridCol{"legch", "", "legsh", ""}, types.GridCol{"legner", "legmh", "corhal", ""}, types.GridCol{"", "", "legah", ""}}, "leggant":types.GridRow{types.GridCol{"corcat", "corkarg", "corshiva", "corkorg"}, types.GridCol{"corjugg", "corsok", "legpede", "leegmech"}, types.GridCol{"legkeres", "", "", ""}}, "leghp":types.GridRow{types.GridCol{"legch", "", "legsh", ""}, types.GridCol{"legner", "legmh", "legcar", ""}, types.GridCol{"", "", "legah", ""}}, "leglab":types.GridRow{types.GridCol{"legck", "cornecro", "leggob", ""}, types.GridCol{"legbal", "leglob", "legkark", "legcen"}, types.GridCol{"", "", "corcrash", ""}}, "legvp":types.GridRow{types.GridCol{"legcv", "legmlv", "leghades", "legscout"}, types.GridCol{"leghelios", "leggat", "legbar", ""}, types.GridCol{"legotter", "legamphtank", "legrail", ""}}}

