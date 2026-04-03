import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
  CardFooter,
} from "@/components/ui/card"

type ProgressCardProps = {
  title: string
  value: number
  max: number
  color: string // tailwind color name, e.g., "blue", "green", etc.
}

export function ProgressCard({
  title,
  value,
  max,
  color,
}: ProgressCardProps) {
const percent = Math.min(100, (value / max) * 100)
  return (
    <div>
        <Card className={`w-full max-w-sm h-40 relative ${color}`}>
            <Card className="w-full max-w-sm h-40 absolute left-1 top-0 right-0">
            <CardHeader>
                <CardDescription>{title}</CardDescription>
                <CardTitle className="text-3xl font-bold">
                {value} <span className="text-gray-400 text-xl">/ {max}</span>
                </CardTitle>
            </CardHeader>
            <CardContent>
                <div className="w-full h-2 bg-gray-200 rounded-full">
                <div
                    className={`h-2 rounded-full ${color}`}
                    style={{ width: `${percent}%` }}
                />
                </div>
            </CardContent>
            <CardFooter>
                {/* action button ถ้าต้องการ */}
            </CardFooter>
            </Card>
        </Card>
    </div>
  )
}
