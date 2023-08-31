import {useMemo} from "react";
import {ScaleLinear} from "d3";

type AxisLeftProps = {
    yScale: ScaleLinear<number, number>;
    pixelsPerTick: number;
    xOffset?: number;
};

// tick length
const TICK_LENGTH = 6;

export const AxisLeft = ({ yScale, pixelsPerTick, xOffset = 20 as number }: AxisLeftProps) => {
    const range = yScale.range();

    const ticks = useMemo(() => {
        const height = range[0] - range[1];
        const numberOfTicksTarget = Math.floor(0.5*height / pixelsPerTick);

        return yScale.ticks(numberOfTicksTarget).map((value) => ({
            value,
            yOffset: yScale(value),
        }));
    }, [yScale, pixelsPerTick]);

    return (
        <g className="axis-left" transform={`translate(${xOffset}, 0)`}>
            {ticks.map((d, i) => (
                <g key={i} transform={`translate(0,${d.yOffset})`}>
                    <text
                        key={i}
                        style={{ fontSize: '24px', textAnchor: 'end', transform: 'translate(-5px, 5px)' }}
                    >
                        {d.value + ' ms'}
                    </text>
                </g>
            ))}
        </g>
    );
};

export default AxisLeft;
