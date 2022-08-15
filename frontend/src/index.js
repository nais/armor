import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import reportWebVitals from './reportWebVitals';
import app from "./App";
import axios from "axios";

const root = ReactDOM.createRoot(document.getElementById('root'));

async function fetchData() {
    try {
        const result = await axios.get("")
        console.log(result.data)
    } catch (error) {
        console.error(error)
    }
}

root.render(fetchData);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
