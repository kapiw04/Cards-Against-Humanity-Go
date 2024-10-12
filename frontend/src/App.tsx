import React from 'react';
import logo from './logo.svg';
import './App.css';

function App() {
    fetch("http://localhost:8080/api/?color=white&k=5", {
      method: "GET", // Or POST, PUT, DELETE, etc.
      headers: {
          "Content-Type": "application/json"
      }
  })
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error("Error fetching data:", error));

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.tsx</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </header>
    </div>
  );
}

export default App;
