// src/components/ui/SuggestiveInput.tsx
import * as React from "react";

import { cn } from "@/lib/utils"; // 调整路径根据你的项目结构
import {
  Command,
  CommandEmpty,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Input } from "@/components/ui/input";

interface SuggestiveInputProps extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'onChange' | 'value'> {
  /** 当前输入框的值 (受控) */
  value: string;
  /** 值改变时的回调函数 */
  onChange: (value: string) => void;
  /** 建议选项列表 */
  suggestions: string[];
  /** 输入框的占位符 */
  placeholder?: string;
  /** 自定义 CSS 类名 */
  className?: string;
  /** Popover 内容的 CSS 类名 */
  popoverClassName?: string;
  /** 没有建议时的提示信息 */
  emptyMessage?: string;
}

const SuggestiveInput = React.forwardRef<HTMLInputElement, SuggestiveInputProps>(
  (
    {
      value,
      onChange,
      suggestions,
      placeholder,
      className,
      popoverClassName,
      emptyMessage = "", // Default empty message
      ...props
    },
    ref
  ) => {
    const [open, setOpen] = React.useState(false);
    const inputRef = React.useRef<HTMLInputElement>(null);
    const [finalSuggestions, setFinalSuggestions] = React.useState(suggestions);

    React.useImperativeHandle(ref, () => inputRef.current as HTMLInputElement);

    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
      const newValue = event.target.value;
      onChange(newValue);
      setFinalSuggestions([newValue, ...suggestions]);
      if (newValue && !open) {
        setOpen(true);
      }
      if (!newValue && open) {
         setOpen(false);
      }
    };

    const handleSuggestionSelect = (selectedValue: string) => {
      if (selectedValue !== value) {
        onChange(selectedValue);
      }
      setOpen(false);
    };

    const filteredSuggestions = React.useMemo(() => {
      const lowerCaseValue = value?.toLowerCase() || "";

      const validSuggestions = finalSuggestions.filter(s => s);

      if (!lowerCaseValue) {
         return validSuggestions;
      }
      return validSuggestions.filter(s =>
        s.toLowerCase().includes(lowerCaseValue)
      );
    }, [value, finalSuggestions]);

    const shouldShowPopover = ((filteredSuggestions.length > 0 || (value && filteredSuggestions.length === 0)) && open) ? true : false;

    return (
      <Popover open={shouldShowPopover}>
        <PopoverTrigger asChild>
          <Input
            ref={inputRef}
            type="text"
            value={value}
            onChange={handleInputChange}
            onFocus={() => {
              // 仅当有建议或已有内容时才在聚焦时打开
              if (finalSuggestions.length > 0) {
                 setOpen(true);
              }
            }}
            placeholder={placeholder}
            className={cn("w-full text-sm", className)}
            role="combobox" // Accessibility role
            aria-expanded={shouldShowPopover} // Accessibility state
            aria-autocomplete="list" // Accessibility hint
            autoComplete="off" // Prevent browser's default autocomplete
            {...props} // Pass down other standard input props like 'id', 'name', etc.
          />
        </PopoverTrigger>
        {shouldShowPopover && (
            <PopoverContent
              className={cn("w-[--radix-popover-trigger-width] p-1", popoverClassName)}
              style={{ zIndex: 50 }}
              onOpenAutoFocus={(e) => e.preventDefault()}
            >
              <Command shouldFilter={false}>
                <CommandList className="pt-0">
                  {value && filteredSuggestions.length === 0 ? (
                    <CommandEmpty>{emptyMessage}</CommandEmpty>
                  ) : null}
                  {filteredSuggestions.map((suggestion) => (
                    <CommandItem
                      key={suggestion}
                      value={suggestion}
                      onSelect={handleSuggestionSelect}
                      className="cursor-pointer"
                    >
                      {suggestion}
                    </CommandItem>
                  ))}
                </CommandList>
              </Command>
            </PopoverContent>
        )}
      </Popover>
    );
  }
);

SuggestiveInput.displayName = "SuggestiveInput";

export { SuggestiveInput };
