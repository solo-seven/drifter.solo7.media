import React from "react";
import { render, screen, waitFor } from '@testing-library/react';
import Home from '../page';

describe('Home page', () => {
  beforeEach(() => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ status: 'ok' })
    }) as jest.Mock;
  });

  test('shows connected status on successful fetch', async () => {
    render(<Home />);
    expect(screen.getByText('Drifter Frontend')).toBeInTheDocument();
    await waitFor(() => expect(screen.getByText(/connected/)).toBeInTheDocument());
  });

  test('shows disconnected when fetch fails', async () => {
    (global.fetch as jest.Mock).mockRejectedValueOnce(new Error('fail'));
    render(<Home />);
    await waitFor(() => expect(screen.getByText('disconnected')).toBeInTheDocument());
  });

  test('shows error text when response is not ok', async () => {
    (global.fetch as jest.Mock).mockResolvedValueOnce({ ok: false, statusText: 'Bad' });
    render(<Home />);
    await waitFor(() => expect(screen.getByText(/error:/i)).toBeInTheDocument());
  });
});
