"use client";

import { useState, useEffect } from 'react';

const DEFAULT_ENV = `{
  "metadata": {
    "name": "Test Environment Alpha",
    "version": "0.1.0",
    "author": "Solo7",
    "created": "2025-06-21T12:00:00Z",
    "tags": ["test", "prototype"]
  }
}`;

export default function Home() {
  const [status, setStatus] = useState('checking...');
  const [envText, setEnvText] = useState(DEFAULT_ENV);
  const [submitStatus, setSubmitStatus] = useState('');

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

  const submitEnvironment = async () => {
    const backendUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080';
    try {
      const response = await fetch(`${backendUrl}/environments`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: envText,
      });
      if (response.ok) {
        setSubmitStatus('submitted');
      } else {
        setSubmitStatus(`error: ${response.statusText}`);
      }
    } catch (e) {
      console.error('Failed to submit environment:', e);
      setSubmitStatus('network error');
    }
  };

  return (
    <main className="flex min-h-screen flex-col items-center justify-start p-8 space-y-4">
      <h1 className="text-4xl font-bold">Drifter Frontend</h1>
      <p>
        Backend status:{' '}
        <span className="font-mono bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 p-1.5 rounded">
          {status}
        </span>
      </p>
      <div className="w-full max-w-2xl">
        <textarea
          className="w-full h-64 p-2 border rounded text-sm dark:bg-gray-800"
          value={envText}
          onChange={(e) => setEnvText(e.target.value)}
        />
        <button
          className="mt-2 px-4 py-2 bg-blue-600 text-white rounded"
          onClick={submitEnvironment}
        >
          Submit Environment
        </button>
        {submitStatus && <p className="mt-2">Status: {submitStatus}</p>}
      </div>
    </main>
  );
}

