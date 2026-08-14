export function mixColors(colorA: string, colorB: string, ratio: number): string {
    const mix = (a: number, b: number, ratio: number): string =>
        Math.round(a + (b - a) * ratio)
            .toString(16)
            .padStart(2, '0')
    const parse = (color: string) => {
        const hex = parseInt(color.slice(1), 16)
        return [(hex >> 16) & 255, (hex >> 8) & 255, hex & 255]
    }
    const [rA, gA, bA] = parse(colorA)
    const [rB, gB, bB] = parse(colorB)
    const red = mix(rA, rB, ratio)
    const green = mix(gA, gB, ratio)
    const blue = mix(bA, bB, ratio)
    return `#${red}${green}${blue}`
}
