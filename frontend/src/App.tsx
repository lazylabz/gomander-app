import {Children, useState} from "react";
import {EventsOn} from "../wailsjs/runtime";
import {ExecCommand} from "../wailsjs/go/main/App";

function App() {
    const [logs, setLogs] = useState<{type: string, line: string}[]>([])
    const [inputText, setInputText] = useState("")

    EventsOn("processLog", (event: {type: string, line: string}) => {
        setLogs((prev) => [...prev, event])
    })

    return (
        <div className="app-container">
            <input value={inputText} onChange={e => setInputText(e.target.value)} />
            <button onClick={() => ExecCommand(inputText || "echo POTATO")}>Click me</button>
            <div>
                {Children.toArray(logs.map((line) => <p>[{line.type}]: {line.line}</p>))}
            </div>
        </div>
    )
}

export default App
