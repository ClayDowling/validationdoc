namespace gameoflife
{
    public class Cell
    {
        public Cell()
        {
            IsAlive = false;
        }
        public Cell(bool isalive)
        {
            IsAlive = isalive;
        }
        
        public bool IsAlive { get; set; }
        public bool Lives(int neighbors)
        {
            if (neighbors == 2 || neighbors == 3) {
                return true;
            }
            return false;
        }
    }
}