import React from 'react';

interface ToggleInputProps {
  label: string;
  checked: boolean;
  onChange: (checked: boolean) => void;
}

export const ToggleInput: React.FC<ToggleInputProps> = ({
  label,
  checked,
  onChange,
}) => {
  return (
    <div className="flex items-center justify-between mb-4">
      <span className="text-sm font-medium text-gray-600">{label}</span>
      <button
        type="button"
        className={`${
          checked ? 'bg-indigo-600' : 'bg-gray-200'
        } relative inline-flex items-center h-6 rounded-full w-11 transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500`}
        onClick={() => onChange(!checked)}
      >
        <span
          className={`${
            checked ? 'translate-x-6' : 'translate-x-1'
          } inline-block w-4 h-4 transform bg-white rounded-full transition-transform`}
        />
      </button>
    </div>
  );
};

export default ToggleInput;
