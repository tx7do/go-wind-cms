import * as React from "react"
import {View} from '@tarojs/components'
import {cn} from "@/lib/utils"

interface CardBaseProps {
  className?: string
  children?: React.ReactNode
  onClick?: (e: any) => void
}

function Card({className, children, ...props}: CardBaseProps) {
  return (
    <View
      className={cn(
        "flex flex-col gap-6 rounded-xl border bg-card py-6 text-card-foreground shadow-sm",
        className
      )}
      {...props}
    >
      {children}
    </View>
  )
}

function CardHeader({className, children, ...props}: CardBaseProps) {
  return (
    <View
      className={cn(
        "grid auto-rows-min grid-rows-[auto_auto] items-start gap-2 px-6",
        className
      )}
      {...props}
    >
      {children}
    </View>
  )
}

function CardTitle({className, children, ...props}: CardBaseProps) {
  return (
    <View
      className={cn("leading-none font-semibold", className)}
      {...props}
    >
      {children}
    </View>
  )
}

function CardDescription({className, children, ...props}: CardBaseProps) {
  return (
    <View
      className={cn("text-sm text-muted-foreground", className)}
      {...props}
    >
      {children}
    </View>
  )
}

function CardAction({className, children, ...props}: CardBaseProps) {
  return (
    <View
      className={cn(
        "col-start-2 row-span-2 row-start-1 self-start justify-self-end",
        className
      )}
      {...props}
    >
      {children}
    </View>
  )
}

function CardContent({className, children, ...props}: CardBaseProps) {
  return (
    <View
      className={cn("px-6", className)}
      {...props}
    >
      {children}
    </View>
  )
}

function CardFooter({className, children, ...props}: CardBaseProps) {
  return (
    <View
      className={cn("flex items-center px-6", className)}
      {...props}
    >
      {children}
    </View>
  )
}

export {
  Card,
  CardHeader,
  CardFooter,
  CardTitle,
  CardAction,
  CardDescription,
  CardContent,
}
