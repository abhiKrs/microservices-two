package sh.logfire.flink.filter
import org.joda.time.{DateTime, DateTimeFieldType, LocalDate, LocalTime}

object JodaConverter {
  private var instance: JodaConverter = null
  private var instantiated = false

  def getConverter: JodaConverter = {
    if (instantiated) return instance
    try {
      Class.forName("org.joda.time.DateTime", false, Thread.currentThread.getContextClassLoader)
      instance = new JodaConverter
    } catch {
      case e: ClassNotFoundException =>
        instance = null
    } finally instantiated = true
    instance
  }
}

class JodaConverter private () {
  def convertDate(`object`: Any): Long = {
    val value = `object`.asInstanceOf[LocalDate]
    value.toDate.getTime
  }

  def convertTime(`object`: Any): Int = {
    val value = `object`.asInstanceOf[LocalTime]
    value.get(DateTimeFieldType.millisOfDay)
  }

  def convertTimestamp(`object`: Any): Long = {
    val value = `object`.asInstanceOf[DateTime]
    value.toDate.getTime
  }
}
