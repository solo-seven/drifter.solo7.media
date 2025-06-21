"use client";

import { useState, useEffect } from 'react';

export default function Home() {
  const [status, setStatus] = useState('checking...');

  useEffect(() => {
    const checkHealth = async () => {
      // The backend URL is expected to be in an environment variable
      const backendUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080';
      try {
        const response = await fetch(`${backendUrl}/health`);
        if (response.ok) {
          const data = await response.json();
          setStatus(`connected (${data.status})`);
        } else {
          setStatus(`error: ${response.statusText}`);
        }
      } catch (error) {
        console.error('Failed to connect to backend:', error);
        setStatus('disconnected');
      }
    };

    checkHealth();
  }, []);

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-24">
      <h1 className="text-4xl font-bold">Drifter Frontend</h1>
      <p className="mt-4">
        Backend status: <span className="font-mono bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 p-1.5 rounded">{status}</span>
      </p>
    </main>
  );
}

