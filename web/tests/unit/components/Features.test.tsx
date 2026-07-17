import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { Features } from '../../../src/components/Features';
import { MATPLOTLIB } from '../../../src/data';

describe('Features', () => {
  it('renders the going-further snippet, the features list and a docs pointer', () => {
    const { container } = render(<Features lib={MATPLOTLIB} />);
    expect(container.querySelector(`#${MATPLOTLIB.id}-more`)).not.toBeNull();
    expect(screen.getByRole('heading', { name: 'Going further' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Features' })).toBeInTheDocument();
    expect(container.querySelectorAll('ul.feat li').length).toBe(MATPLOTLIB.features.length);
    const docs = screen.getByRole('link', { name: /docs tab/ });
    expect(docs).toHaveAttribute('href', '#docs');
  });
});
