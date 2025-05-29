'use client'

import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useRef, useState } from "react";

type UserJoined = {
    username: string
}

type Info = {
    selfUsername: string
    opponentUsername: string
}

function DisplayNames(
    { username, opponent } 
    :{ username : string, opponent : string }
) {
    let msg;
    if(username === '' && opponent === '') {
        msg = "Loading...";
    } else if(opponent === '') {
        msg = username + " waiting for opponent...";
    } else {
        msg = username + " vs " + opponent;
    }

    return (
        <div className="flex flex-row justify-center">
            <h1 className="text-2xl">{msg}</h1>
        </div>
    )
}

type Connect4Board = number[][];
const Columns = 7;
const Rows = 6;

function Connect4({ board} : { board: Connect4Board }) {
    const [hoveredColumn, setHoveredColumn] = useState<number | null>(null);

    function createRow() {
        const arr = [];
        for(let i = 0; i < Rows; i++) {
            arr.push(0);
        }
        return arr;
    }
    
    while(board.length < Columns) {
        board.push(createRow());
    }

    function handleColumnClick(columnIndex: number) {
        alert(`Clicked column: ${columnIndex}`);
    };
    
    return (
        <div className="inline-block bg-blue-600 p-2 rounded-lg shadow-lg">
            <div className="flex flex-row gap-1">
                { 
                    board.map((col, i) => 
                        <div 
                            key={i} 
                            className={`flex flex-col gap-1 cursor-pointer transition-all duration-200 p-1 rounded ${
                                hoveredColumn === i ? 'bg-blue-500 transform scale-105' : ''
                            }`}
                            onMouseEnter={() => setHoveredColumn(i)}
                            onMouseLeave={() => setHoveredColumn(null)}
                            onClick={() => handleColumnClick(i)}
                        >
                            { 
                                col.map((space, j) => {
                                    let bgColor = "bg-white";
                                    if (space < 0) {
                                        bgColor = "bg-red-500";
                                    } else if (space > 0) {
                                        bgColor = "bg-yellow-400";
                                    }
                                    
                                    return (
                                        <div 
                                            key={j} 
                                            className={`w-14 h-14 ${bgColor} rounded-full border-2 border-blue-800 shadow-inner`}
                                        />
                                    );
                                })
                            }
                        </div>
                    )
                }
            </div>
        </div>
    );
}

export default function Page() {
    const router = useRouter()
    const key = useSearchParams().get("key");
    const [username, setUsername] = useState('');
    const [opponent, setOppenent] = useState('');
    const wsRef = useRef<WebSocket | null>(null);

    useEffect(() => {
        if(key == null) {
            router.push('/');
            return;
        }

        let uri = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        uri += '//localhost:8080/game?key=' + key;
        const ws = new WebSocket(uri);

        ws.onmessage = (e) => {
            const obj = JSON.parse(e.data);

            const info = obj as Info
            if(info.selfUsername != null && info.opponentUsername != null) {
                setUsername(info.selfUsername);
                setOppenent(info.opponentUsername);
                return;
            }

            if((obj as UserJoined).username != null) {
                setOppenent((obj as UserJoined).username);
                return;
            }

            console.log(obj);
            alert("json not recognized");
        }

        return () => {
            ws.close()
        }
    }, [key]);

    useEffect(() => {
        console.log("user:", username, "opp:",opponent)
    }, [username, opponent])

    return (
        <div className="min-h-screen flex items-center justify-center">
            <div className="flex flex-col gap-10 items-center">
                <DisplayNames username={username} opponent={opponent} />
                <Connect4 board={[]} />
            </div>
        </div>
    );
}