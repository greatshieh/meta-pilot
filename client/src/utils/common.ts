/**
 * Transform record to option
 *
 * @example
 *   ```ts
 *   const record = {
 *     key1: 'label1',
 *     key2: 'label2'
 *   };
 *   const options = transformRecordToOption(record);
 *   // [
 *   //   { value: 'key1', label: 'label1' },
 *   //   { value: 'key2', label: 'label2' }
 *   // ]
 *   ```;
 *
 * @param record
 */
export function transformRecordToOption<T extends Record<string, string>>(record: T) {
    return Object.entries(record).map(([value, label]) => ({
        value,
        label
    })) as CommonType.Option<keyof T>[]
}

export function transformArrayToOption<T>(array: T[]) {
    return array.map(item => ({
        value: item,
        label: item
    })) as CommonType.Option[]
}

/**
 * Translate options
 *
 * @param options
 */
export function translateOptions(options: CommonType.Option<string>[]) {
    return options.map(option => ({
        ...option,
        label: option.label
    }))
}

/**
 * Toggle html class
 *
 * @param className
 */
export function toggleHtmlClass(className: string) {
    function add() {
        document.documentElement.classList.add(className)
    }

    function remove() {
        document.documentElement.classList.remove(className)
    }

    return {
        add,
        remove
    }
}
