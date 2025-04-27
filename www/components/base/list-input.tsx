import React, { useState } from 'react';
import { Badge } from '../ui/badge';
import { Input } from '../ui/input';
import { Button } from '../ui/button';
import { useTranslation } from 'react-i18next';
import { cn } from '@/lib/utils';

interface StringListInputProps {
  value: string[];
  onChange: React.Dispatch<React.SetStateAction<string[]>>;
  placeholder?: string;
  className?: string;
}

const StringListInput: React.FC<StringListInputProps> = ({ value, onChange, placeholder, className }) => {
  const { t } = useTranslation();
  const [inputValue, setInputValue] = useState('');

  const handleAdd = () => {
    if (inputValue.trim()) {
      if (value && value.includes(inputValue)) {
        return;
      }

      if (value) {
        onChange([...value, inputValue]);
      } else {
        onChange([inputValue]);
      }
      setInputValue('');
    }
  };

  const handleRemove = (itemToRemove: string) => {
    onChange(value.filter(item => item !== itemToRemove));
  };

  return (
    <div className={cn("mx-auto", className)}>
      <div className="flex items-center mb-4">
        <Input
          type="text"
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          className="flex-1 px-4 py-2 border rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
          placeholder={placeholder || t('input.list.placeholder')}
        />
        <Button
          disabled={!inputValue || value && value.includes(inputValue)}
          onClick={handleAdd}
          className="ml-2 px-4 py-2"
        >
          {t('input.list.add')}
        </Button>
      </div>
      {
        value && <div className="flex flex-wrap gap-2">
          {value.map((item, index) => (
            <Badge key={index} className='flex flex-row items-center justify-start'>{item}
              <div
                onClick={() => handleRemove(item)}
                className="ml-1 h-4 w-4 text-center rounded-full hover:text-red-500 cursor-pointer"
              >
                ×
              </div>
            </Badge>
          ))}
        </div>
      }
    </div>
  );
};

export default StringListInput;