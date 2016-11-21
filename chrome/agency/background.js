// Copyright (c) 2011 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Called when the user clicks on the browser action.
chrome.browserAction.onClicked.addListener(function(tab) {
  // No tabs or host permissions needed!
  console.log('Turning ' + JSON.stringify({name:"John Rambo", time:"2pm"}) + ' red!');
  var xmlhttp = new XMLHttpRequest();   // new HttpRequest instance 
  xmlhttp.open("POST", "http://localhost:9090/");
  xmlhttp.send(tab.url);
  chrome.tabs.executeScript({
    code: 'document.body.style.backgroundColor="red"'
  });
});


chrome.tabs.onActivated.addListener(function(activeInfo){
chrome.tabs.get(activeInfo.tabId, function(tab){

  var xmlhttp = new XMLHttpRequest();   // new HttpRequest instance 
  xmlhttp.open("POST", "http://localhost:9090/");
  xmlhttp.send(JSON.stringify(tab));
})

})

// chrome.tabs.onUpdated.addListener(function(tabId, changeInfo, tab) {
    
//   var xmlhttp = new XMLHttpRequest();   // new HttpRequest instance 
//   xmlhttp.open("POST", "http://localhost:9090/");
//   xmlhttp.send(tab.url);
// });

chrome.tabs.onReplaced.addListener(function(activeInfo){
chrome.tabs.get(activeInfo.tabId, function(tab){

  var xmlhttp = new XMLHttpRequest();   // new HttpRequest instance 
  xmlhttp.open("POST", "http://localhost:9090/");
  xmlhttp.send(JSON.stringify(tab));
})

})


