import type { Meta, StoryObj } from '@storybook/react';
import { useState } from 'react';
import { ToggleInput } from '../toggle-input';

const meta: Meta<typeof ToggleInput> = {
  title: 'Components/UI/ToggleInput',
  component: ToggleInput,
  tags: ['autodocs'],
  argTypes: {
    label: { control: 'text' },
    checked: { control: 'boolean' },
  },
  args: {
    label: 'Toggle Setting',
    checked: false,
  },
};

export default meta;
type Story = StoryObj<typeof ToggleInput>;

// Interactive component that handles state
const ToggleWithState = (args: any) => {
  const [checked, setChecked] = useState(args.checked || false);
  return <ToggleInput {...args} checked={checked} onChange={setChecked} />;
};

export const Default: Story = {
  render: (args) => <ToggleWithState {...args} />,
};

export const CheckedByDefault: Story = {
  render: (args) => <ToggleWithState {...args} checked={true} />,
  args: {
    label: 'Dark Mode',
  },
};

export const WithDifferentLabels: Story = {
  render: (args) => <ToggleWithState {...args} />,
  args: {
    label: 'Enable Notifications',
  },
};

// Multiple toggles in a single story
export const MultipleToggles: Story = {
  render: () => (
    <div className="space-y-4">
      <ToggleWithState label="Auto-save changes" checked={true} />
      <ToggleWithState label="Email notifications" checked={false} />
      <ToggleWithState label="Dark theme" checked={true} />
      <ToggleWithState label="Two-factor authentication" checked={false} />
    </div>
  ),
};
