/**
 * @author Naorem Khogendro Singh
 */
var app = angular.module('app');

app.controller('auth_ctrl', function($scope, $http) {
	$scope.init = function() {
	};
	$scope.login = function(googleUser) {
		const userData = JSON.stringify(googleUser)
		console.log("Logged in succesfully");
		console.log("User data: " + userData)
		const profile = googleUser.getBasicProfile();
		const auth = googleUser.getAuthResponse();
		const user = profile.getName();
		const email = profile.getEmail();
		const idToken = auth["id_token"];
		$scope.session = idToken;
		setCookie("session", $scope.session);
		setCookie("user", user);
		setCookie("email", email);
		redirectOnLogin();
	};
});

app.controller('controlboard_ctrl', function($scope, $http, $timeout, $window) {
	$scope.init = function() {
		$scope.power_state = 0
		$scope.session = getCookie("session");
		if ($scope.session === null) {
			$scope.logout();
		}
		$scope.user = getCookie("user");
		$scope.email = getCookie("email");

		const request = createHttpRequest("GET", `/v1/controllers/${controllerId}`, null, $scope.session)
		console.log('Calling getControlBoard');
		$http(request).then(function(response) {
			console.log('getControlBoard: '+JSON.stringify(response.data));
			if (response.data == null) {
				return;
			}
			$scope.periodicFunction();
			$scope.intervalFunction();
			$scope.initDisplay();
		}).catch(function (err) {
			const data = err.data;
			console.log("Error: " + JSON.stringify(data));
			$scope.alert(data.message);
			if (data.code == 401) {
				$scope.logout();
			}
		});
	};

	$scope.logout = function() {
		redirectOnLogout();
	};

	$scope.alert = function(msg) {
		alertFunc(msg);
	}

	$scope.periodicFunction = function() {
		$scope.listDevices(function() {
			console.log("Done fetching devices");
		});
	};

	$scope.listDevices = function(callback) {
		const uri = `/v1/controllers/${controllerId}`
		const request = createHttpRequest("GET", uri, null, $scope.session)
		$http(request).then(function(response) {
			console.log('getController: ' + response.data);
			const controller = response.data
			console.log("Controller: " + JSON.stringify(controller))
			$scope.devices = {};
			for (var idx in controller.pins) {
				const pin = controller.pins[idx]
				console.log("Processing pin: " + JSON.stringify(pin))
				const device = $scope.displayDevice(pin);
				$scope.devices[String(device.deviceId)] = device
			}
			callback();
		}).catch(function (err) {
			const data = err.data;
			console.log("Error: " + JSON.stringify(data));
			$scope.alert(data.message);
			if (data.code == 401) {
				$scope.logout();
			}
		});
	};

	$scope.isAnyDeviceAvailable = function(){
		if ($scope.devices) {
			return Object.keys($scope.devices).length > 0;
		}
		return false;
	};

	// Function to replicate setInterval using $timeout service
	// (5s).
	$scope.intervalFunction = function() {
		$timeout(function() {
			var e = $scope.periodicFunction();
			if (!checkException(e)) {
				$scope.intervalFunction();
			}
		}, funcInterval)
	};

	$scope.initDisplay = function() {
		$scope.hours = [];
	    for (h = 0; h < 24; h++) {
	       $scope.hours.push(h);
	    }
	    $scope.mins = [];
	    for (m = 0; m < 60; m++) {
	       $scope.mins.push(m);
	    }
	    $scope.statuss = ["ON", "OFF"];
	};

	$scope.displayDevice = function(pin) {
		var extendedDevice = null;
		if (pin.value == "High") {
			extendedDevice = {deviceId : pin.id, name : pin.name, powerOn : true, button_css_class: poweron_css_class};
		} else {
			extendedDevice = {deviceId : pin.id, name : pin.name, powerOn : false, button_css_class: poweroff_css_class};
		}
		console.log("Found device: " + JSON.stringify(extendedDevice))
		return extendedDevice;
	};

	/* Power click handler starts */
	$scope.powerHandler = function(device) {
		var powerValue = "High";
		if (device.powerOn) {
			powerValue = "Low";
		}
		const deviceId = device.deviceId
		const uri = `/v1/controllers/${controllerId}`
		const payload = JSON.parse(`{
			"id": "${controllerId}",
			"pins": [{
				"id": ${deviceId},
				"mode": "Output",
				"value": "${powerValue}"
			}]
		}`)
		const request = createHttpRequest("PUT", uri, payload, $scope.session)
		$http(request).then(function(response) {
			console.log("Response for updateController: " + response.data);
			if (response.data == null) {
				return;
			}
			const controller = response.data
			for (var idx in controller.pins) {
				const pin = controller.pins[idx]
				const device = $scope.displayDevice(pin);
				$scope.devices[String(device.deviceId)] = device
			}
		}).catch(function (err) {
			const data = err.data;
			console.log("Error: " + JSON.stringify(data));
			$scope.alert(data.message);
			if (data.code == 401) {
				$scope.logout();
			}
		});
	};
});
