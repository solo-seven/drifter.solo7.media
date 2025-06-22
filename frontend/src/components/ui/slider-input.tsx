import React from 'react';

interface SliderInputProps {
  label: string;
  value: number;
  onChange: (value: number) => void;
  min: number;
  max: number;
  step?: number;
  unit?: string;
}

export const SliderInput: React.FC<SliderInputProps> = ({
  label,
  value,
  onChange,
  min,
  max,
  step = 1,
  unit,
}) => {
  return (
    <div className="mb-4">
      <label className="block text-sm font-medium text-gray-600 mb-1">{label}</label>
      <div className="flex items-center">
        <input
          type="range"
          min={min}
          max={max}
          step={step}
          value={value}
          onChange={(e) => onChange(parseFloat(e.target.value))}
          className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer"
        />
        <span className="ml-4 text-sm text-gray-800 w-16 text-right">
          {value}
          {unit}
        </span>
      </div>
    </div>
  );
};

export default SliderInput;
