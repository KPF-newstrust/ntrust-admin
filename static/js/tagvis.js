
var tagExplains = {		
	"NNG":	"일반 명사",
	"NNP":	"고유 명사",
	"NNB":	"의존 명사",
	"NNBC":	"단위를 나타내는 명사",
	"NR":	"수사",
	"NP":	"대명사",
	"VV":	"동사",
	"VA":	"형용사",
	"VX":	"보조 용언",
	"VCP":	"긍정 지정사",
	"VCN":	"부정 지정사",
	"MM":	"관형사",
	"MAG":	"일반 부사",
	"MAJ":	"접속 부사",
	"IC":	"감탄사",
	"JKS":	"주격 조사",
	"JKC":	"보격 조사",
	"JKG":	"관형격 조사",
	"JKO":	"목적격 조사",
	"JKB":	"부사격 조사",
	"JKV":	"호격 조사",
	"JKQ":	"인용격 조사",
	"JX":	"보조사",
	"JC":	"접속 조사",
	"EP":	"선어말 어미",
	"EF":	"종결 어미",
	"EC":	"연결 어미",
	"ETN":	"명사형 전성 어미",
	"ETM":	"관형형 전성 어미 ",
	"XPN":	"체언 접두사",
	"XSN":	"명사 파생 접미사",
	"XSV":	"동사 파생 접미사",
	"XSA":	"형용사 파생 접미사",
	"XR":	"어근",
	"SF":	"마침표, 물음표, 느낌표",
	"SE":	"줄임표 …",
	"SSO":	"여는 괄호 (, [",
	"SSC":	"닫는 괄호 ), ]",
	"SC":	"구분자 , · / :",
	"SY": "	",
	"SL":	"외국어",
	"SH":	"한자",
	"SN":	"숫자",
	"?": "UNKNOWN",
	"SS_": "문장끝",

	"NNM":"단위 의존 명사",
	"VXV":"보조 동사",
	"VXA":"보조 형용사",
	"MDT":"일반 관형사",
	"MDN":"수 관형사",
	"MAC":"접속 부사",
	"JKM":"부사격 조사",
	"JKI":"호격 조사",
	"EPH":"존칭 선어말 어미",
	"EPT":"시제 선어말 어미",
	"EPP":"공손 선어말 어미",
	"EFN":"평서형 종결 어미",
	"EFQ":"의문형 종결 어미",
	"EFO":"명령형 종결 어미",
	"EFA":"청유형 종결 어미",
	"EFI":"감탄형 종결 어미",
	"EFR":"존칭형 종결 어미",
	"ECE":"대등 연결 어미",
	"ECD":"의존적 연결 어미",
	"ECS":"보조적 연결 어미",
	"ETD":"관형형 전성 어미",
	"XPV":"용언 접두사",
	"SP":"쉼표,가운뎃점,콜론,빗금",
	"SS":"따옴표,괄호,줄표",
	"SO":"붙임표(물결,숨김,빠짐)",
	"SW":"기타기호 (논리수학기호,화폐기호)",
	"UN":"분석안됨, 명사로 추정됨",
	"OL":"외국어",
	"OH":"한자",
	"ON":"숫자",

	"N":"명사,체언",
	"J":"조사,관계언",
	"P":"용언",
	"X":"접사",
	"M":"수식언",
	"E":"어미",
	"I":"독립언",
	"S":"기호",
	"F":"외국어",
	
	"Noun":"명사",
	"Verb":"동사",
	"Adverb":"부사",
	"Adjective":"형용사",
	"Josa":"조사",
	"Eomi":"어미",
	"PreEomi":"선어말 어미",
	"Suffix":"접미사",
	"Number":"숫자",
	"Alpha":"외국어",
	"Email":"이메일",
	"Foreign":"외래어",
	"Determiner":"한정사",
	"Conjunction":"접속사",
	"Punctuation":"문장부호"
};


var balloon = null;
var CY_LINE = 52;
var XPOS_WORDBEGIN = 40;

function VisualizePostag(selector, tagdata, maxSvgWidth)
{
	if (!balloon)
	{
		balloon = $('<div class="arrow_box"></div>');
		$("body").append(balloon);
	}

	if (!maxSvgWidth)
		maxSvgWidth = 500;

	var svgElem = $(selector);
	var svg = d3.select(selector);

	var curX = 0;
	var curY = 0;
	var grpCurRow;
	var cxMax = maxSvgWidth;

	var bgColorToggle = true;
	function openNewRow()
	{
		if (curX > 0)
		{
			curX = 0;
			curY += CY_LINE;
		}

		svg.append("rect")
			.attr("x", 0)
			.attr("y", curY)
			.attr("width", "100%")
			.attr("height", CY_LINE)
			.attr("fill", (bgColorToggle = !bgColorToggle) ? "#f8f8f8":"#eee");

		grpCurRow = svg.append("g");
	}

	function addPos(word, tagAll)
	{
		curX += 5;
		var grp = grpCurRow.append("g");
		grp.attr("data-xx","123");

		var tags = tagAll.split(/\+|,/)
		var colorClass = "tag_" + tags[0];
		var tipHtml = "";
		for (var t in tags)
		{
			if (tipHtml) tipHtml += "<br>+";
			tipHtml += '<b class="tag_'+tags[t]+'">';
			if (tags[t] == "UNKNOWN") tags[t] = "?";
			var exp = tagExplains[tags[t]] || "설명추가요망";
			tipHtml += tags[t]+"</b>:"+exp;
		}
		var tag0 = tags[0];
		
		var cxTag = tag0.length * 8;
		var cxWord = (word.length || 1) * 15;
		var xCenter = curX + cxWord/2;

		grp.append("rect")
			.attr("class", colorClass)
			.attr("x", xCenter - cxTag/2)
			.attr("y", curY + 4)
			.attr("width", cxTag)
			.attr("height", 14)
			.attr("rx", 3)
			.attr("ry", 3);
		
		grp.append("text")
			.attr("x", xCenter)
			.attr("y", curY + 14)
			.attr("text-anchor", "middle")
			.attr("font-size", "10px")
			.text(tag0);

		var yArc = curY + 27;
		grp.append("path")
			.attr("d", "M" + (curX) + "," + (yArc)
				+ "C" + (curX + 1) + "," + (yArc - 5)
				+ " " + (curX + cxWord - 1) + "," + (yArc - 5)
				+ " " + (curX + cxWord) + "," + (yArc))
			.attr("stroke", "#888")
			.attr("stroke-width", 1.5)
			.attr("fill", "none");

		var textBox = grp.append("rect")
			.attr("x", curX)
			.attr("y", curY + 30)
			.attr("width", cxWord + "px")
			.attr("height", 18)
			.attr("fill", "rgba(0,0,0,0.05)");
		
		grp.append("text")
			.attr("y", curY + 44)
			.attr("text-anchor", "middle")
			.attr("font-size", "14px")
			.attr("x", xCenter)
			.text(word);

		grp.on("mouseover", function() {
			textBox.attr("stroke","#000");
			balloon.html(tipHtml);
			balloon.show();

			var rcSvg = svgElem[0].getBoundingClientRect();
			var rcElem = this.getBBox();
			var x = rcSvg.left + xCenter - balloon.width() /2 -8;
			var y = rcSvg.top + rcElem.y + window.scrollY - balloon.height() - 12;
			balloon.offset({left:x, top:y});
		});

		grp.on("mouseout", function() {
			textBox.attr("stroke","none");
			balloon.hide();
		});

		curX += cxWord;
		if (curX > maxSvgWidth)
		{
			if (curX > cxMax)
				cxMax = curX;

			openNewRow();
		}
	}

	function drawHR(y)
	{
		svg.append("line")
			.attr("x1", 0)
			.attr("y1", y)
			.attr("x2", "100%")
			.attr("y2", y)
			.attr("stroke-width", 1)
			.attr("stroke","#000")
	}

	for (var iSentence in tagdata)
	{
		var sentence = tagdata[iSentence];

		openNewRow();
		drawHR(curY+1);		

		// draw sentence number with background oval
		var grpSentNum = grpCurRow.append("g");
		grpSentNum.append("rect")
			.attr("x", 3)
			.attr("y", curY + 16)
			.attr("width", 28)
			.attr("height", 24)
			.attr("rx", 5)
			.attr("ry", 5)
			.attr("stroke","#000")
			.attr("fill", "#ffd")
		grpSentNum.append("text")
			.attr("x", 3+14)
			.attr("y", curY + 32)
			.attr("font-size", "14px")
			.attr("text-anchor", "middle")
			.text(iSentence);

		curX = XPOS_WORDBEGIN;

		for (var iWord in sentence)
		{
			var item = sentence[iWord];
			addPos(item[0],item[1]);
		}		
	}

	if (curX > 0) curY += CY_LINE;
	drawHR(curY-1);
	svg.attr("height", curY).attr("width", cxMax);
}