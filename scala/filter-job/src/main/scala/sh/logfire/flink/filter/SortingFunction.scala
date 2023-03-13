package sh.logfire.flink.filter

import org.apache.flink.api.common.state.{ MapState, MapStateDescriptor }
import org.apache.flink.configuration.Configuration
import org.apache.flink.streaming.api.functions.KeyedProcessFunction
import org.apache.flink.util.Collector

import java.util

class SortingFunction extends KeyedProcessFunction[java.lang.Integer, String, String] {

  var state: MapState[java.lang.Long, util.ArrayList[String]] = null
  private val waitTimeMs = 1000

  override def open(parameters: Configuration): Unit =
    state = getRuntimeContext.getMapState(
      new MapStateDescriptor[java.lang.Long, util.ArrayList[String]]("sorting",
                                                                     classOf[java.lang.Long],
                                                                     classOf[util.ArrayList[String]])
    )

  override def processElement(i: String,
                              context: KeyedProcessFunction[Integer, String, String]#Context,
                              collector: Collector[String]): Unit = {

    val eventTime = context.timestamp()

    val list = if (state.contains(eventTime)) {
      state.get(eventTime)
    } else {
      new util.ArrayList[String]()
    }
    list.add(i)

    state.put(context.timestamp(), list)
//  TODO
    context.timerService.registerEventTimeTimer(context.timestamp() + waitTimeMs)

  }

  override def onTimer(
      timestamp: Long,
      ctx: KeyedProcessFunction[Integer, String, String]#OnTimerContext,
      out: Collector[String]
  ): Unit = {

    val stateValue = state.get(timestamp - waitTimeMs)
    if (stateValue != null) {
      stateValue.forEach(r => out.collect(r))
//      TODO check if flink executes processElement in another thread and creates the key
      state.remove(timestamp - waitTimeMs)
    }

  }
}
