var app = angular.module('app', []);

const domain = "controller.bytescheme.com"
const baseUrl = `https://${domain}`
const controllerId = "bfd8dd0a-10db-4782-86ec-b27f52d6362c"
const clientId = "91456297737-d1p2ha4n2847bpsrdrcp72uhp614ar9q.apps.googleusercontent.com"
const poweron_css_class = 'power_on'
const poweroff_css_class = 'power_off'
const funcInterval = 10000
//const poweroff_img = 

app.directive('parseStyle', function($interpolate) {
	return function(scope, elem) {
		var exp = $interpolate(elem.html()), watchFunc = function() {
			return exp(scope);
		};
		scope.$watch(watchFunc, function(html) {
			elem.html(html);
		});
	};
});

function guid() {
	function s4() {
		return Math.floor((1 + Math.random()) * 0x10000).toString(16)
				.substring(1);
	}
	return s4() + s4() + '-' + s4() + '-' + s4() + '-' + s4() + '-' + s4()
			+ s4() + s4();
};

function getTriggerTimeSec(hr, min) {
	var seconds = Math.round(new Date().getTime() / 1000);
	seconds += (hr * 60 * 60);
	seconds += (min * 60);
	return seconds;
};

function createHttpRequest(method, uri, payload, session_id) {
	var request = {
		method : method,
		url : baseUrl + uri,
		headers : {
			'Content-Type' : 'application/json',
			'Authorization' : session_id
		},
		data : payload
	}
	return request
}

function extractOrigin(url) {
	var start = url.indexOf("://");
	var end = url.indexOf("/", start + 3);
	return "https" + url.substring(start, end);
};

function redirectOnLogin() {
	const controlboardPage = baseUrl + "/controlboard.html";
	/*
	const url = "https://accounts.google.com/o/oauth2/v2/auth?scope=email&client_id="
	        + clientId
			+ "&redirect_uri="
			+ controlboardPage
			+ "&response_type=token";
	*/
	document.location.replace(controlboardPage);
};

function redirectOnLogout() {
	const logoutUrl = `https://www.google.com/accounts/Logout?continue=https://appengine.google.com/_ah/logout?continue=${baseUrl}`;
	document.location.replace(logoutUrl);
};

function alertFunc(msg) {
	alert(msg);
};

function doInScope(appName, callback) {
	var appElement = document.querySelector('[ng-app=' + appName + ']');
	var $scope = angular.element(appElement).scope();
	$scope.$apply(callback($scope));
};

function onSignIn(googleUser) {
	console.log("Signed in with user: " + JSON.stringify(googleUser))
	doInScope('app', function(scope) {
		scope.login(googleUser);
	});
};

function signOut() {
	var auth2 = gapi.auth2.getAuthInstance();
	auth2.signOut().then(function() {
		doInScope('app', function(scope) {
			scope.logout();
		});
	});
};

function checkException(e) {
	if (e) {
		var msg = JSON.stringify(e);
		if (msg.includes("Session") || msg.includes("security check")) {
			redirectOnLogout();
			return true;
		}
	}
	return false;
};

function setCookie(cname, cvalue) {
	var mins = 10;
	var d = new Date();
	d.setTime(d.getTime() + (mins * 60 * 1000));
	var expires = "expires=" + d.toUTCString();
	document.cookie = cname + "=" + cvalue + ";" + expires
			+ ";domain=" + domain;
};

function getCookie(cname) {
	var name = cname + "=";
	var ca = document.cookie.split(';');
	for (var i = 0; i < ca.length; i++) {
		var c = ca[i];
		while (c.charAt(0) == ' ') {
			c = c.substring(1);
		}
		if (c.indexOf(name) == 0) {
			return c.substring(name.length, c.length);
		}
	}
	return null;
};

function deleteCookie(cname) {
	var d = new Date();
	d.setTime(d.getTime() - (60 * 1000));
	var expires = "expires=" + d.toUTCString();
	document.cookie = cname + "=;" + expires;
};