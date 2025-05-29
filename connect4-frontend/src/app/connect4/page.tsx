'use client'

import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";

type UserJoined = {
    username: string
}

export default function Page() {
    const router = useRouter()
    const key = useSearchParams().get("key");
    const [opponent, setOppenent] = useState('');

    useEffect(() => {
        if(key == null) {
            router.push('/');
            return;
        }

        let uri = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        uri += '//localhost:8080/game?key=' + key;
        const ws = new WebSocket(uri);

        ws.onopen = () => {
        }

        ws.onmessage = (e) => {
            const obj = JSON.parse(e.data);

            if((obj as UserJoined).username != null) {
                setOppenent((obj as UserJoined).username);
                return;
            }


            alert("json not recognized");
        }

        ws.onclose = (e) => {

        }

        ws.onerror = () => {

        }
    }, [key]);

    return (
        <div>
            {
                opponent.length == 0 && 
                <h1>{key} ready to play!</h1>
            }

            {
                opponent.length != 0 &&
                <h1>{key} vs {opponent}</h1>
            }
        </div>
    );
}