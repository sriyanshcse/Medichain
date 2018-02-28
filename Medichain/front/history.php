<?php

$medid = $_GET['id'];
?>



<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="description" content="">
    <meta name="keywords" content="">
    <meta name="author" content="">

    <title>
      
        Ledger 
      
    </title>

    <link href="https://fonts.googleapis.com/css?family=Roboto:300,400,500,700,400italic" rel="stylesheet">
    
      <link href="style1.css" rel="stylesheet">
    
    
    <link href="style2.css" rel="stylesheet">

    <style>
      /* note: this is a hack for ios iframe for bootstrap themes shopify page */
      /* this chunk of css is not part of the toolkit :) */
      body {
        width: 1px;
        min-width: 100%;
        *width: 100%;
      }

        html, body {
  height: 100%;
  margin: 0;
  padding: 0;
}
    </style>

      <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js"></script>
    <script src="https://ajax.googleapis.com/ajax/libs/angularjs/1.6.4/angular.min.js"></script>
    <script async defer src="https://maps.googleapis.com/maps/api/js?key=AIzaSyB9ZgDbEH6lp_X6YE4JNt_M8dc3UyQs6j0">
</script>


      <script>
      var app = angular.module('myApp', []);
      app.controller('myCtrl', function($scope, $http, $window) {
          $http({
            method : "GET",
            headers:{  
                        "authorization": "Bearer " + $window.localStorage["org1token"],
                     "Accept":"application/json",
                      "content-type":"application/x-www-form-urlencoded"
                    },   
            url:"http://localhost:4000/channels/mychannel/chaincodes/fabcar?peer=peer1&fcn=queryMedHistory&args=%5B%22" + '<?php echo $medid ?>' + "%22%5D",
          }).then(function mySuccess(response) {
              
              console.log(response.data);

                     // var r=JSON.parse(response.data[0]);

                $scope.myData = response.data;
                $scope.initMap(response.data);
            })

          var poly;
          var map;

          $scope.initMap = function(data) {
            console.log("hello");
            var points = [];
            for (var i = 0; i < data.length; i++) {
                var loc = data[i].location.split(" ");
                points.push({
                  lat: loc[0],
                  lng: loc[1]
                });
            }
            map = new google.maps.Map(document.getElementById('map'), {
              zoom:3,
              center: {lat: parseFloat(points[0].lat), lng: parseFloat(points[0].lng)}  // Center the map on Chicago, USA.
            });
            poly = new google.maps.Polyline({
              strokeColor: '#000000',
              strokeOpacity: 1.0,
              strokeWeight: 3
            });
            poly.setMap(map);

            // Add a listener for the click event
            for (var i = 0; i < points.length; i++) {
               addLatLng(new google.maps.LatLng(points[i].lat,points[i].lng))
            }
          }

console.log("hello world");

          // Handles click events on a map, and adds a new point to the Polyline.
          function addLatLng(point) {
            var path = poly.getPath();

            // Because path is an MVCArray, we can simply append a new coordinate
            // and it will automatically appear.
            path.push(point);
            
            // Add a new marker at the new plotted point on the polyline.
            var marker = new google.maps.Marker({
              position: point,
              title: '#' + path.getLength(),
              map: map
            });
          }
        
      });
      </script>

      <script type="text/javascript">
        
        


      </script>

  </head>


<body ng-app="myApp" ng-controller="myCtrl">
  <div class="bw">
    <div class="dh">
      <div class="et brf">
        <div class="bqn">
  <div class="bqo">
    <br/><br/><br/>
    <h2 class="bqp">History</h2>
  </div>
</div>

<div class="">
  <table class="ck" data-sort="table">

    <thead>
      <col width="10">
      <tr>

        <th>TransactionID</th>
        <th>Time</th>
        <th>Temperature</th>
        <th>Location</th>
        <th>Owner</th>
        <th>Status</th>
      </tr>
    </thead>
    <tbody>
      <tr ng-repeat="x in myData">
        <td>{{x.txid}}</td>
        <td>{{x.timestamp}}</td>
        <td>{{x.temperature}} C</td>
        <td>{{x.location}}</td>
        <td>{{x.owner}}</td>
        <td ng-if="x.status" class="Success"><b>SUCCESS</b></td>
        <td ng-if="!x.status" class="Failure"><b>FAILED</b></td>
      </tr>
    </tbody>
  </table>
</div>


      </div>
    </div>
  </div>

<center>
  
  <div id="map" style="height: 500px; width: 800px; "></div>


</center>


    <!-- <script src="script1.js"></script>
    <script src="script2.js"></script>
    <script src="script3.js"></script>
    <script src="script4.js"></script>
    <script src="script5.js"></script>
    <script src="script6.js"></script>
 -->    <script>
      // execute/clear BS loaders for docs
      $(function(){while(window.BS&&window.BS.loader&&window.BS.loader.length){(window.BS.loader.pop())()}})
    </script>
  </body>
</html>


