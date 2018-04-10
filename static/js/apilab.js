function waitUntilDone(objId,limit, csrf)
{
	if (--limit <= 0)
	{
		hideWaitIndi();
		alert("작업을 완료하는데 오래 걸릴 것 같습니다.\n목록으로 돌아갑니다.");
		window.location.href = "list";
		return;
	}

	setTimeout(function() {

		$.ajax({
			method: "POST",
			url: "/admin/lab/query.json",
			data: {
				"cmd": "wait",
				"_csrf": csrf,
				"id": objId
			},
			success: function(resp) {
				if (resp.code)
				{
					hideWaitIndi();
					alert(resp.msg);
					return;
				}

				if (resp.data.done)
				{
					window.location.href = window.location.pathname + "?id=" + objId;
				}
				else
				{
					waitUntilDone(objId, limit, csrf);
				}
			},
			error: function (xhr, ajaxOptions, thrownError) {
				hideWaitIndi();
				alert("AJAX Error: " + thrownError);
			}
		});

	}, 1000);
}

function postLabNew(payload)
{
	if (!confirm("전송하시겠습니까?"))
		return;

	showWaitIndi();
	$.ajax({
		method: "POST",
		url: "/admin/lab/new.json",
		data: payload,
		success: function(resp) {
			if (resp.code)
			{
				hideWaitIndi();
				alert(resp.msg);
				return;
			}

			waitUntilDone(resp.data.id, 5, payload._csrf);
		},
		error: function (xhr, ajaxOptions, thrownError) {
			hideWaitIndi();
			alert("AJAX Error: " + thrownError);
		}
	});
}

function onLabResend(objId,csrf)
{
	if (!confirm("API 호출을 다시 요청하시겠습니까?"))
		return;

	showWaitIndi();
	$.ajax({
		method: "POST",
		url: "/admin/lab/query.json",
		data: {
			"cmd": "resend",
			"_csrf": csrf,
			"id": objId
		},
		success: function(resp) {
			if (resp.code)
			{
				hideWaitIndi();
				alert(resp.msg);
				return;
			}

			waitUntilDone(objId, 5, csrf);
		},
		error: function (xhr, ajaxOptions, thrownError) {
			hideWaitIndi();
			alert("AJAX Error: " + thrownError);
		}
	});
}

function onLabDelete(objId, csrf)
{
	if (!confirm("삭제하시겠습니까?"))
		return;

	showWaitIndi();
	$.ajax({
		method: "POST",
		url: "/admin/lab/query.json",
		data: {
			"cmd": "delete",
			"_csrf": csrf,
			"id": objId
		},
		success: function(resp) {
			if (resp.code)
			{
				alert(resp.msg);
				return;
			}

			alert("삭제되었습니다.");
			window.location.href = "list";
		},
		error: function (xhr, ajaxOptions, thrownError) {
			alert("AJAX Error: " + thrownError);
		},
		complete: function() {
			hideWaitIndi();
		}
	});
}
