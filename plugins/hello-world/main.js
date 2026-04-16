const readline = require('readline');

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout
});

let idCounter = 0;

// 发送请求到主进程
function sendRequest(method, params) {
  const id = `req-${idCounter++}`;
  const request = {
    id: id,
    method: method,
    payload: params
  };
  console.log(JSON.stringify(request));
  return new Promise((resolve) => {
    // 在实际插件中，这里会通过事件监听器匹配 ID 来接收响应
    // 为了演示简单，我们假设主进程会立即返回
  });
}

// 激活插件
async function activate() {
  console.error("Hello World Plugin is now active!");
  
  // 调用 Hao-Code API 显示消息
  await sendRequest("vscode.window.showInformationMessage", "Hello from Hao-Code Sandbox!");
}

activate();
