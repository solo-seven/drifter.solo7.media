import type { Meta, StoryObj } from '@storybook/react';
import { CollapsibleSection } from '../collapsible-section';

const meta: Meta<typeof CollapsibleSection> = {
  title: 'Components/UI/CollapsibleSection',
  component: CollapsibleSection,
  tags: ['autodocs'],
  argTypes: {
    title: { control: 'text' },
    defaultOpen: { control: 'boolean' },
  },
  args: {
    title: 'Section Title',
    defaultOpen: true,
  },
};

export default meta;
type Story = StoryObj<typeof CollapsibleSection>;

export const Default: Story = {
  args: {
    children: (
      <div className="p-4">
        <p>This is the collapsible content area. You can put any content here.</p>
        <p>It can include multiple paragraphs, components, or any other React elements.</p>
      </div>
    ),
  },
};

export const ClosedByDefault: Story = {
  args: {
    ...Default.args,
    defaultOpen: false,
  },
};

export const WithLongContent: Story = {
  args: {
    title: 'Long Content Section',
    children: (
      <div className="space-y-2">
        {Array.from({ length: 10 }).map((_, i) => (
          <div key={i} className="p-2 bg-gray-100 rounded">
            Item {i + 1}
          </div>
        ))}
      </div>
    ),
  },
};
