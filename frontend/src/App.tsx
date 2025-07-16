import {Children, useEffect, useState} from "react";
import {EventsOff, EventsOn} from "../wailsjs/runtime";
import {ExecCommand} from "../wailsjs/go/main/App";

function App() {
    const [logs, setLogs] = useState<{ type: string, line: string }[]>([])
    const [inputText, setInputText] = useState("")

    useEffect(() => {

        EventsOn("newLog", (event: { type: string, line: string }) => {
            setLogs((prev) => [...prev, event])
        })

        return () => {
            EventsOff("newLog")
        }
    })


    return (
        <div className="app-container">
            <p>{logs.length}</p>
            <input value={inputText} onChange={e => setInputText(e.target.value)}/>
            <button onClick={() => ExecCommand(inputText || "echo POTATO")}>Click me</button>
            <button onClick={() => setLogs([])}>Clean</button>
            <div>
                {Children.toArray(logs.map((line) => <p>[{line.type}]: {line.line}</p>))}
            </div>
        </div>
    )
}

export default App
