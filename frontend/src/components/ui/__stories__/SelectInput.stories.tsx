import type { Meta, StoryObj } from '@storybook/react';
import { useState } from 'react';
import { SelectInput } from '../select-input';

const meta: Meta<typeof SelectInput> = {
  title: 'Components/UI/SelectInput',
  component: SelectInput,
  tags: ['autodocs'],
  argTypes: {
    label: { control: 'text' },
    options: { control: 'object' },
  },
  args: {
    label: 'Select an option',
    options: [
      { value: 'option1', label: 'Option 1' },
      { value: 'option2', label: 'Option 2' },
      { value: 'option3', label: 'Option 3' },
    ],
  },
};

export default meta;
type Story = StoryObj<typeof SelectInput>;

// Interactive component that handles state
const SelectWithState = (args: any) => {
  const [value, setValue] = useState(args.options[0].value);
  return <SelectInput {...args} value={value} onChange={setValue} />;
};

export const Default: Story = {
  render: (args) => <SelectWithState {...args} />,
};

export const WithManyOptions: Story = {
  render: (args) => <SelectWithState {...args} />,
  args: {
    label: 'Choose a country',
    options: [
      { value: 'us', label: 'United States' },
      { value: 'ca', label: 'Canada' },
      { value: 'mx', label: 'Mexico' },
      { value: 'br', label: 'Brazil' },
      { value: 'uk', label: 'United Kingdom' },
      { value: 'de', label: 'Germany' },
      { value: 'fr', label: 'France' },
      { value: 'es', label: 'Spain' },
      { value: 'it', label: 'Italy' },
      { value: 'jp', label: 'Japan' },
    ],
  },
};

export const WithGroups: Story = {
  render: (args) => <SelectWithState {...args} />,
  args: {
    label: 'Select a food category',
    options: [
      { value: '', label: 'Select a category', disabled: true },
      { value: 'fruits', label: 'Fruits' },
      { value: 'apples', label: 'Apples', group: 'fruits' },
      { value: 'bananas', label: 'Bananas', group: 'fruits' },
      { value: 'oranges', label: 'Oranges', group: 'fruits' },
      { value: 'vegetables', label: 'Vegetables' },
      { value: 'carrots', label: 'Carrots', group: 'vegetables' },
      { value: 'broccoli', label: 'Broccoli', group: 'vegetables' },
      { value: 'dairy', label: 'Dairy' },
      { value: 'milk', label: 'Milk', group: 'dairy' },
      { value: 'cheese', label: 'Cheese', group: 'dairy' },
    ],
  },
};
