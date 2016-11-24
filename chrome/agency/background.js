function getXMLHTTP() {
  var xmlhttp = new XMLHttpRequest();   // new HttpRequest instance 
  xmlhttp.open("POST", "http://localhost:9090/");
  return xmlhttp
}

chrome.tabs.onActivated.addListener(function (activeInfo) {
  chrome.tabs.get(activeInfo.tabId, function (tab) {
    getXMLHTTP().send(tab.url);
  })

})

chrome.tabs.onUpdated.addListener(function (tabId, changeInfo, tab) {
  if (changeInfo.url != null) {
    getXMLHTTP().send(changeInfo.url);
  }
});


chrome.windows.onFocusChanged.addListener(function (windowId) {
  var xmlhttp = getXMLHTTP()
  if (windowId === chrome.windows.WINDOW_ID_NONE) {
    xmlhttp.send("losss all focus");
  } else {
    chrome.tabs.query({ active: true, windowId: windowId }, function (tabs) {
      xmlhttp.send(tabs[0].url);
    })
  }
})


