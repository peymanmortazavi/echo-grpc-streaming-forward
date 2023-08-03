function handleClick() {
	let responseContainer = document.getElementById("response");
	let countInput = document.getElementById("count");

	responseContainer.innerHTML = "";

	let socket = new WebSocket("ws://localhost:5000/ws");

	socket.addEventListener("open", e => {
		socket.send(`{"count": ${countInput.value}}`);
	});

	socket.addEventListener("message", (e) => {
		n = document.createElement("p");
		n.innerText = "data: " + e.data;
		responseContainer.append(n);
	});
}
