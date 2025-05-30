'use client'

import { useRouter, useSearchParams } from "next/navigation";
import { Suspense, useCallback, useEffect, useRef, useState } from "react";

type UserJoined = {
    username: string
}

type Info = {
    selfUsername: string
    opponentUsername: string
    gameWinner: number
    starter: boolean
    yourTurn : boolean
    plays: number[]
}

type LogPlay = {
    column: number
    isSelf: boolean
}

type ErrorMessage = {
    error: string
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

function Connect4(
    { ws, board, yourTurn } 
    : { ws : WebSocket | null, board: Connect4Board, yourTurn : boolean | null }
) {
    const [hoveredColumn, setHoveredColumn] = useState<number | null>(null);

    function handleColumnClick(columnIndex: number) {
        if(ws == null || ws.readyState !== WebSocket.OPEN) {
            return;
        }

        const message = {
            column: columnIndex
        };

        ws.send(JSON.stringify(message));
    }
    
    return (
        <div>
            {
                yourTurn !== null &&
                <h1 className={
                    `${yourTurn ? "text-green-700" : "text-red-500"} text-center text-xl mb-2`
                }>
                    { yourTurn ? 'Your' : "Opponent's"} turn
                </h1>
            }
            <div className="inline-block bg-blue-600 p-2 rounded-lg shadow-lg">
                <div className="flex flex-row gap-1">
                    { 
                        board.map((col, i) => 
                            <div 
                                key={i} 
                                className={`flex flex-col gap-5 cursor-pointer transition-all duration-200 p-1 rounded ${
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
                                                className={`w-8 h-8 ${bgColor} rounded-full border-2 border-blue-800 shadow-inner`}
                                            />
                                        );
                                    })
                                }
                            </div>
                        )
                    }
                </div>
            </div>
        </div>
    );
}

function initBoard() : Connect4Board {
    const board : Connect4Board = [];
    function makeRow() {
        const arr = [];
        for(let i = 0; i < Rows; i++) {
            arr.push(0);
        }

        return arr;
    }

    while(board.length < Columns) {
        board.push(makeRow());
    }

    return board;
}

function FindNewGame(
    {ws, gameWinner, func} 
    : {ws : WebSocket | null, gameWinner : string | null, func : () => void}
) {
    function handleClick() {
        if(ws == null || ws.readyState !== WebSocket.OPEN) {
            return;
        }

        ws.send(JSON.stringify({ requested: true }));
        func();
    }

    if(!gameWinner) {
        return <></>
    }

    return (
        <div className="flex flex-row gap-4 items-center">
            <h1 className={`${
                    gameWinner === 'Opponent' 
                    ? 'text-red-500'
                    : 'text-green-700'
                } text-2xl`}
            >
                {gameWinner} win{gameWinner === 'Opponent' ? 's' : ''}!
            </h1> 
            <button onClick={handleClick} type="button" className="cursor-pointer text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-xl px-5 py-2.5 me-2 mb-2 dark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none dark:focus:ring-blue-800">Play again</button>
        </div>
    )
}

function GameContent() {
    const router = useRouter()
    const key = useSearchParams().get("key");
    const [username, setUsername] = useState('');
    const [opponent, setOpponent] = useState('');
    const [connectionStatus, setConnectionStatus] = useState('Connecting...');
    const wsRef = useRef<WebSocket | null>(null);
    const [board, setBoard] = useState<Connect4Board>(initBoard());
    const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
    const [gameWinner, setGameWinner] = useState<string | null>(null);
    const [boardUpToDate, setBoardUpToDate] = useState<boolean>(false);
    const [yourTurn, setYourTurn] = useState<boolean | null>(null);

    const reset = useCallback(() => {
        setOpponent('');
        setBoard(initBoard());
        setBoardUpToDate(false);
        setGameWinner(null);
        setYourTurn(null);
    }, [
        setOpponent, 
        setBoard,
        setGameWinner,
        setBoardUpToDate,
        setYourTurn
    ])

    const connectWebSocket = useCallback(() => {
        if(key == null) {
            router.push('/');
            return;
        }

        // Close existing connection if any
        if (wsRef.current) {
            wsRef.current.close();
        }

        let uri = window.location.origin.split(":")[1];
        uri = (window.location.protocol === 'https:' ? 'wss:' : 'ws:') + uri;
        uri += ':8080/game?key=' + key;
        
        const ws = new WebSocket(uri);
        wsRef.current = ws;

        ws.onopen = () => {
            setConnectionStatus('Connected');
        };

        ws.onclose = (event) => {
            setConnectionStatus('Disconnected');
            
            // Don't auto-reconnect on normal closure or if component is unmounting
            if (event.code !== 1000 && event.code !== 1001) {
                // Reconnect after 3 seconds
                reconnectTimeoutRef.current = setTimeout(() => {
                    connectWebSocket();
                }, 3000);
            }
        };

        ws.onerror = (error) => {
            console.error("WebSocket error:", error);
            setConnectionStatus('Error');
        };

        ws.onmessage = (e) => {
            
            try {
                const obj = JSON.parse(e.data);

                // Handle initial player info
                const info = obj as Info;
                if (
                    info.selfUsername !== undefined 
                    && info.opponentUsername !== undefined
                    && info.gameWinner !== undefined
                    && info.starter !== undefined
                    && info.plays !== undefined
                    && info.yourTurn !== undefined
                ) {
                    switch(info.gameWinner) {
                        case 1:
                            setGameWinner('You');
                            break;
                        case 2:
                            setGameWinner('Opponent')
                            break;
                        default:
                            setGameWinner(null);
                            break;
                    }

                    setUsername(info.selfUsername);
                    setOpponent(info.opponentUsername);
                    setYourTurn(info.opponentUsername === '' || info.gameWinner !== 0 ? null : info.yourTurn);

                    if(info.plays.length === 0 || boardUpToDate) {
                        return;
                    }

                    function makeEmptyBoard() {
                        const b : Connect4Board = [];
                        for(let i = 0; i < Columns; i++) {
                            const arr = [];
                            for(let j = 0; j < Rows; j++) {
                                arr.push(0);
                            }

                            b.push(arr);
                        }

                        return b;
                    }

                    const b = makeEmptyBoard();
                    function dropPiece(b : Connect4Board, c : number, player : number) {
                        for(let r = Rows - 1; r >= 0; r--) {
                            if(b[c][r] === 0) {
                                b[c][r] = player;
                                return;
                            }
                        }
                    }

                    let p = info.starter;
                    for(let i = 0; i < info.plays.length; i++) {
                        // fill with players
                        dropPiece(b, info.plays[i], p ? 1 : -1);
                        p = !p;
                    }

                    // add empty spots
                    for(let i = 0; i < Columns; i++) {
                        while(b[i].length < Rows) {
                            b[i].push(0);
                        }
                    }

                    setBoardUpToDate(true);
                    setBoard(b);

                    return;
                } else if(info.selfUsername !== undefined) {
                    setUsername(info.selfUsername);
                }

                // Handle player joined
                if (obj.username !== undefined) {
                    const userJoined = obj as UserJoined;
                    setOpponent(userJoined.username);
                    return;
                }

                // Handle game moves
                if (obj.column !== undefined && obj.isSelf !== undefined) {
                    const logPlay = obj as LogPlay;
                    // Update board with the move
                    const newBoard = [...board];
                    const column = newBoard[logPlay.column];
                    
                    // Find the lowest empty space in the column
                    for (let row = Rows - 1; row >= 0; row--) {
                        if (column[row] === 0) {
                            column[row] = logPlay.isSelf ? 1 : -1;
                            break;
                        }
                    }
                    
                    setBoard(newBoard);
                    return;
                }

                // Handle game over
                if (obj.youWin !== undefined) {
                    setGameWinner(obj.youWin ? 'You' : 'Opponent');
                    return
                }

                // Handle errors
                if (obj.error !== undefined) {
                    const error = obj as ErrorMessage;
                    if(error.error === "failed to do action") {
                        alert("Not your turn yet");
                    }
                    return;
                }

                // Handle other message types
                alert("Message type not handled");
            } catch (error) {
                console.error("Error parsing message:", error, "Raw message:", e.data);
            }
        };
    }, [board, setBoard, setUsername, setOpponent, setGameWinner, boardUpToDate, setBoardUpToDate, board]);

    useEffect(() => {
        connectWebSocket();

        return () => {
            // Cleanup on unmount
            if (reconnectTimeoutRef.current) {
                clearTimeout(reconnectTimeoutRef.current);
            }
            if (wsRef.current) {
                wsRef.current.close();
            }
        };
    }, [key, connectWebSocket]);

    return (
        <div className="min-h-screen flex items-center justify-center">
            <div className="flex flex-col gap-10 items-center">
                <div className="text-sm text-gray-500">Status: {connectionStatus}</div>
                <DisplayNames username={username} opponent={opponent} />
                <Connect4 yourTurn={yourTurn} board={board} ws={wsRef.current} />
                <FindNewGame func={reset} ws={wsRef.current} gameWinner={gameWinner} />
            </div>
        </div>
    );
}

export default function Page() {
    return (
        <Suspense>
            <GameContent />
        </Suspense>
    )
}