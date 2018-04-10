function showWaitIndi()
{
	$('#waitIndi').show();
}
function hideWaitIndi()
{
	$('#waitIndi').hide();
}

function applyDtFilters(sel,dt)
{
	var dirty=false;
	$("input,select",sel).each(function(idx,elem) {
		if (elem.value && elem.dataset.column)
		{
			dirty = true;
			dt.columns(elem.dataset.column).search(elem.value);
		}
	});
	
	if (dirty)
	{
		dt.ajax.reload();
	}
}

var entityMap = {
	'&': '&amp;',
	'<': '&lt;',
	'>': '&gt;',
	'"': '&quot;',
	"'": '&#39;',
	'/': '&#x2F;',
	'`': '&#x60;',
	'=': '&#x3D;'
};
  
function escapeHtml(string)
{
	return String(string).replace(/[&<>"'`=\/]/g, function (s) { return entityMap[s]; });
}

// NProgress
if (typeof NProgress != 'undefined') {
	$(document).ready(function () {
		NProgress.start();
	});

	$(window).load(function () {
		NProgress.done();
	});
}
