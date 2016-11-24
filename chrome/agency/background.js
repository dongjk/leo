chrome.tabs.onActivated.addListener(function (activeInfo) {
  chrome.tabs.get(activeInfo.tabId, function (tab) {

    var xmlhttp = new XMLHttpRequest();   // new HttpRequest instance 
    xmlhttp.open("POST", "http://localhost:9090/");
    xmlhttp.send(tab.url);
  })

})

chrome.tabs.onUpdated.addListener(function(tabId, changeInfo, tab) {

  var xmlhttp = new XMLHttpRequest();   // new HttpRequest instance 
  xmlhttp.open("POST", "http://localhost:9090/");
  if(changeInfo.url!=null){
    xmlhttp.send(changeInfo.url);
  }
});


chrome.windows.onFocusChanged.addListener(function (windowId) {
  var xmlhttp = new XMLHttpRequest();   // new HttpRequest instance 
  xmlhttp.open("POST", "http://localhost:9090/");
  if (windowId === chrome.windows.WINDOW_ID_NONE) {
    xmlhttp.send("losss all focus");
  } else {
    chrome.tabs.query({active: true, windowId:windowId}, function (tabs){
        xmlhttp.send(tabs[0].url);
    })
  }
})


