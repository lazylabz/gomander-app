import {useEffect, useState} from "react";
import {EventsOn} from "../wailsjs/runtime";
import {GetCommands} from "../wailsjs/go/main/LogServer";
import type {Command} from "./types/contracts.ts";
import {Event} from "./types/contracts.ts";

function App() {
    const [commands, setCommands] = useState<Record<string, Command>>({})

    const refreshCommands = async () => {
        const commandsData = await GetCommands();

        setCommands(commandsData);
    }

    useEffect(() => {
        EventsOn(Event.GET_COMMANDS, () => {
            refreshCommands();
        })
    })

    return (
        <div className="w-full h-full bg-white">
            Hello World! {JSON.stringify(commands)}
        </div>
    )
}

export default App
