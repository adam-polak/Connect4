'use client'

import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";

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
            alert(e.data)
        }

        ws.onclose = (e) => {

        }

        ws.onerror = () => {

        }
    }, [key]);

    return (
        <div>
            <h1>{key} Ready to play!</h1>
        </div>
    );
}