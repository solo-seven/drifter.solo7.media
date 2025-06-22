import type { Meta, StoryObj } from '@storybook/react';
import { useState } from 'react';
import { SliderInput } from '../slider-input';

const meta: Meta<typeof SliderInput> = {
  title: 'Components/UI/SliderInput',
  component: SliderInput,
  tags: ['autodocs'],
  argTypes: {
    label: { control: 'text' },
    min: { control: 'number' },
    max: { control: 'number' },
    step: { control: 'number' },
    unit: { control: 'text' },
  },
  args: {
    label: 'Slider Label',
    min: 0,
    max: 100,
    step: 1,
    value: 50,
    unit: '%',
  },
};

export default meta;
type Story = StoryObj<typeof SliderInput>;

// Interactive component that handles state
const SliderWithState = (args: any) => {
  const [value, setValue] = useState(args.value || 50);
  return <SliderInput {...args} value={value} onChange={setValue} />;
};

export const Default: Story = {
  render: (args) => <SliderWithState {...args} />,
};

export const WithCustomRange: Story = {
  render: (args) => <SliderWithState {...args} min={10} max={1000} step={10} />,
  args: {
    label: 'Custom Range',
    value: 100,
    unit: 'px',
  },
};

export const WithDecimalSteps: Story = {
  render: (args) => <SliderWithState {...args} min={0} max={1} step={0.1} />,
  args: {
    label: 'Opacity',
    value: 0.5,
    unit: '',
  },
};

export const WithoutUnit: Story = {
  render: (args) => <SliderWithState {...args} unit="" />,
  args: {
    label: 'Rating',
    value: 3,
    min: 1,
    max: 5,
  },
};
