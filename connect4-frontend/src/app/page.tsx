'use client'

import { nanoid } from "nanoid";
import { useCallback, useState } from "react";
import { useRouter } from 'next/navigation';

export default function Home() {
  const minUsernameLength = 5;
  const maxUsernameLength = 15;
  const [input, setInput] = useState<string>('');
  const [errorMessage, setErrorMessage] = useState<string>('');
  const router = useRouter();

  const joinGame = useCallback(async () => {
    if(input.length < minUsernameLength) {
      setErrorMessage('username is too short');
      return;
    } else if(input.length > maxUsernameLength) {
      setErrorMessage('username is too long');
      return;
    }

    setErrorMessage('');

    try {
      const key = nanoid();
      const result = await fetch("http://localhost:8080/player/create", {
        method: "POST",
        body: JSON.stringify(
          { 
            LoginKey: key,
            Username: input
          }
        )  
      });

      if(result.status === 200) {
        router.push('/connect4?key=' + key);
      }
    } catch {}
  }, [input]);

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="bg-white p-8 rounded-2xl shadow-lg w-full max-w-sm">
        <h2 className="text-2xl font-semibold text-center mb-6">Join Connect4</h2>
        <div className="space-y-4">
          <input
            onChange={e => setInput(e.target.value)}
            type="text" 
            placeholder="Enter your username" 
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          { 
            errorMessage.length !== 0 &&
            <p className="text-md text-red-500">* {errorMessage}</p>
          }
          <button
            onClick={joinGame}
            type="submit" 
            className="w-full bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 transition duration-200 cursor-pointer"
          >
            Join
          </button>
        </div>
      </div>
    </div>
  );
}
