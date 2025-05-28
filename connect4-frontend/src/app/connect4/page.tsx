'use client'

import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";

export default function Page() {
    const router = useRouter()
    const params = useSearchParams();
    const [key, setKey] = useState<string | null>('');
    useEffect(() => {
        if(key !== '') {
            return;
        }

        const k = params.get('key');
        if(k == null) {
            router.push('/');
            return;
        }

        let uri = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        uri += '//localhost:8080/game?key=' + k;
        const ws = new WebSocket(uri);

        ws.onopen = () => {
            alert('connection opened!');
        }

        ws.onmessage = (e) => {

        }

        ws.onclose = (e) => {

        }

        ws.onerror = () => {

        }
    }, [key, setKey]);

    return (
        <div>
            <h1>{key} Ready to play!</h1>
        </div>
    );
}