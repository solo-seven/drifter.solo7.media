import React from "react";
import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Home from '../page';

// Mock the global fetch
const mockFetch = jest.fn();
global.fetch = mockFetch;

describe('Home page', () => {
  // Helper to wait for the next tick
test('renders without crashing', () => {
  render(<Home />);
});

const waitForNextTick = () => new Promise(resolve => setTimeout(resolve, 0));

  beforeEach(() => {
    mockFetch.mockClear();
    // Default mock for health check
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ status: 'ok' })
    });
  });

  describe('Health check', () => {
    test('shows connected status on successful fetch', async () => {
      render(<Home />);
      expect(screen.getByText('Drifter Frontend')).toBeInTheDocument();
      await waitFor(() => expect(screen.getByText(/connected/)).toBeInTheDocument());
    });

    test('shows disconnected when fetch fails', async () => {
      mockFetch.mockReset();
      mockFetch.mockRejectedValueOnce(new Error('fail'));
      render(<Home />);
      await waitFor(() => expect(screen.getByText('disconnected')).toBeInTheDocument());
    });

    test('shows error text when response is not ok', async () => {
      mockFetch.mockReset();
      mockFetch.mockResolvedValueOnce({ 
        ok: false, 
        statusText: 'Bad',
        json: async () => ({ error: 'Bad request' })
      });
      render(<Home />);
      await waitFor(() => expect(screen.getByText(/error: Bad/i)).toBeInTheDocument());
    });
  });

  describe('Environment submission', () => {
    const testEnv = {
      metadata: {
        name: 'Test Env',
        version: '1.0.0',
        author: 'Tester'
      }
    };
    const testEnvString = JSON.stringify(testEnv, null, 2);

    test('allows typing in the textarea', async () => {
      render(<Home />);
      const textarea = screen.getByRole('textbox');
      await userEvent.clear(textarea);
      fireEvent.change(textarea, { target: { value: testEnvString } });
      expect(textarea).toHaveValue(testEnvString);
    });

    test('submits environment successfully', async () => {
      // Mock the form submission response
      mockFetch.mockResolvedValueOnce({
        ok: true,
        statusText: 'OK'
      });
      
      render(<Home />);
      
      // Wait for initial health check to complete
      await waitFor(() => expect(screen.getByText(/connected/)).toBeInTheDocument());
      
      // Change the environment
      const textarea = screen.getByRole('textbox');
      fireEvent.change(textarea, { target: { value: testEnvString } });
      
      // Submit the form
      const submitButton = screen.getByRole('button', { name: /submit environment/i });
      await userEvent.click(submitButton);
      
      // Verify the fetch was called with the right parameters
      expect(mockFetch).toHaveBeenLastCalledWith(
        expect.stringMatching(/\/environments$/),
        expect.objectContaining({
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: testEnvString,
        })
      );
      
      // Verify success message
      await waitFor(() => {
        const statusElement = screen.getByText(/Status:/);
        expect(statusElement).toHaveTextContent('submitted');
      });
    });

    test('handles submission error', async () => {
      const errorMessage = 'Invalid JSON';
      // Mock the form submission error response
      mockFetch.mockResolvedValueOnce({ 
        ok: false, 
        statusText: errorMessage
      });
      
      render(<Home />);
      
      // Wait for initial health check to complete
      await waitFor(() => expect(screen.getByText(/connected/)).toBeInTheDocument());
      
      const submitButton = screen.getByRole('button', { name: /submit environment/i });
      await userEvent.click(submitButton);
      
      await waitFor(() => {
        const statusElement = screen.getByText(/Status:/);
        expect(statusElement).toHaveTextContent(`error: ${errorMessage}`);
      });
    });

    test('handles network error on submission', async () => {
      const errorMessage = 'Network error';
      // Mock the network error on form submission
      mockFetch.mockRejectedValueOnce(new Error(errorMessage));
      
      render(<Home />);
      
      // Wait for initial health check to complete
      await waitFor(() => expect(screen.getByText(/connected/)).toBeInTheDocument());
      
      const submitButton = screen.getByRole('button', { name: /submit environment/i });
      await userEvent.click(submitButton);
      
      await waitFor(() => {
        const statusElement = screen.getByText(/Status:/);
        expect(statusElement).toHaveTextContent('network error');
      });
    });
  });
});
